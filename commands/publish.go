package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/pivotal-cf/go-pivnet/v2"
	"github.com/pivotal-cf/go-pivnet/v2/logshim"
	"github.com/pivotal-cf/jhanda"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/yaml.v2"
)

const oslFileType = "Open Source License"

//go:generate counterfeiter -o ./fakes/pivnet_releases_service.go --fake-name PivnetReleasesService . PivnetReleasesService
type PivnetReleasesService interface {
	List(productSlug string) ([]pivnet.Release, error)
	Update(productSlug string, release pivnet.Release) (pivnet.Release, error)
}

//go:generate counterfeiter -o ./fakes/pivnet_product_files_service.go --fake-name PivnetProductFilesService . PivnetProductFilesService
type PivnetProductFilesService interface {
	List(productSlug string) ([]pivnet.ProductFile, error)
	AddToRelease(productSlug string, releaseID int, productFileID int) error
}

type Publish struct {
	Options struct {
		Kilnfile    string `short:"f" long:"file" default:"Kilnfile" description:"path to Kilnfile"`
		Version     string `short:"v" long:"version-file" default:"version" description:"path to version file"`
		PivnetToken string `short:"t" long:"pivnet-token" description:"pivnet refresh token" required:"true"`
		PivnetHost  string `long:"pivnet-host" default:"https://network.pivotal.io" description:"pivnet host"`
	}

	PivnetReleaseService      PivnetReleasesService
	PivnetProductFilesService PivnetProductFilesService

	FS  billy.Filesystem
	Now func() time.Time

	OutLogger, ErrLogger *log.Logger
}

func NewPublish(outLogger, errLogger *log.Logger, fs billy.Filesystem) Publish {
	return Publish{
		OutLogger: outLogger,
		ErrLogger: errLogger,
		FS:        fs,
	}
}

func (p Publish) Execute(args []string) error {
	defer p.recoverFromPanic()

	kilnfile, buildVersion, err := p.parseArgsAndSetup(args)
	if err != nil {
		return err
	}

	return p.updateReleaseOnPivnet(kilnfile, buildVersion)
}

func (p Publish) recoverFromPanic() func() {
	return func() {
		if r := recover(); r != nil {
			p.ErrLogger.Println(r)
			os.Exit(1)
		}
	}
}

func (p *Publish) parseArgsAndSetup(args []string) (Kilnfile, *semver.Version, error) {
	_, err := jhanda.Parse(&p.Options, args)
	if err != nil {
		return Kilnfile{}, nil, err
	}

	if p.Now == nil {
		p.Now = time.Now
	}

	if p.PivnetReleaseService == nil || p.PivnetProductFilesService == nil {
		config := pivnet.ClientConfig{
			Host:      p.Options.PivnetHost,
			UserAgent: "kiln",
		}

		tokenService := pivnet.NewAccessTokenOrLegacyToken(p.Options.PivnetToken, p.Options.PivnetHost, false)

		logger := logshim.NewLogShim(p.OutLogger, p.ErrLogger, false)
		client := pivnet.NewClient(tokenService, config, logger)

		if p.PivnetReleaseService == nil {
			p.PivnetReleaseService = client.Releases
		}

		if p.PivnetProductFilesService == nil {
			p.PivnetProductFilesService = client.ProductFiles
		}
	}

	versionFile, err := p.FS.Open(p.Options.Version)
	if err != nil {
		return Kilnfile{}, nil, err
	}
	defer versionFile.Close()

	versionBuf, err := ioutil.ReadAll(versionFile)
	if err != nil {
		return Kilnfile{}, nil, err
	}

	version, err := semver.NewVersion(strings.TrimSpace(string(versionBuf)))
	if err != nil {
		return Kilnfile{}, nil, err
	}

	file, err := p.FS.Open(p.Options.Kilnfile)
	if err != nil {
		return Kilnfile{}, nil, err
	}
	defer file.Close()

	var kilnfile Kilnfile
	if err := yaml.NewDecoder(file).Decode(&kilnfile); err != nil {
		return Kilnfile{}, nil, fmt.Errorf("could not parse Kilnfile: %s", err)
	}

	return kilnfile, version, nil
}

func (p Publish) updateReleaseOnPivnet(kilnfile Kilnfile, buildVersion *semver.Version) error {
	p.OutLogger.Printf("Requesting list of releases for %s", kilnfile.Slug)

	window, err := kilnfile.ReleaseWindow(p.Now())
	if err != nil {
		return err
	}

	rv, err := ReleaseVersionFromBuildVersion(buildVersion, window)
	if err != nil {
		return err
	}

	releaseType := releaseType(window, rv)

	var releases releaseSet
	releases, err = p.PivnetReleaseService.List(kilnfile.Slug)
	if err != nil {
		return err
	}

	release, err := releases.Find(buildVersion.String())
	if err != nil {
		return err
	}

	versionToPublish, err := p.determineVersion(releases, rv)
	if err != nil {
		return err
	}

	err = p.attachLicenseFile(window, kilnfile.Slug, release.ID, versionToPublish)
	if err != nil {
		return err
	}

	release.Version = versionToPublish.String()
	release.ReleaseType = releaseType

	release.ReleaseDate = p.Now().Format(PublishDateFormat)
	if rv.IsGA() {
		lastPatchRelease, matchExists, err := p.findLatestPatchRelease(releases, rv)
		if err != nil {
			return err
		}
		if matchExists {
			release.EndOfSupportDate = lastPatchRelease.EndOfSupportDate
		} else {
			release.EndOfSupportDate = endOfSupportFor(p.Now())
		}
	}

	if _, err := p.PivnetReleaseService.Update(kilnfile.Slug, release); err != nil {
		return err
	}

	return nil
}

func endOfSupportFor(publishDate time.Time) string {
	monthWithOverflow := publishDate.Month() + 10
	month := ((monthWithOverflow - 1) % 12) + 1
	yearDelta := int((monthWithOverflow - 1) / 12)
	startOfTenthMonth := time.Date(publishDate.Year()+yearDelta, month, 1, 0, 0, 0, 0, publishDate.Location())
	endOfNinthMonth := startOfTenthMonth.Add(-24 * time.Hour)
	return endOfNinthMonth.Format(PublishDateFormat)
}

func (p Publish) findLatestPatchRelease(releases releaseSet, rv *releaseVersion) (pivnet.Release, bool, error) {
	constraint, err := rv.MajorMinorConstraint()
	if err != nil {
		return pivnet.Release{}, false, err
	}
	return releases.FindLatest(constraint)
}

func (p Publish) attachLicenseFile(window, slug string, releaseID int, version *semver.Version) error {
	if window == "ga" {
		productFiles, err := p.PivnetProductFilesService.List(slug)
		if err != nil {
			return err
		}

		licenseFileVersion := fmt.Sprintf("%d.%d", version.Major(), version.Minor())

		var productFileID int
		for _, file := range productFiles {
			if file.FileType == oslFileType && file.FileVersion == licenseFileVersion {
				productFileID = file.ID
				break
			}
		}

		if productFileID == 0 {
			return errors.New("required license file doesn't exist on Pivnet")
		}

		err = p.PivnetProductFilesService.AddToRelease(slug, releaseID, productFileID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p Publish) determineVersion(releases releaseSet, version *releaseVersion) (*semver.Version, error) {
	if version.IsGA() {
		return version.Semver(), nil
	}

	constraint, err := version.PrereleaseVersionsConstraint()
	if err != nil {
		return nil, fmt.Errorf("determineVersion: error building prerelease version constraint: %w", err)
	}

	latestRelease, previousReleaseExists, err := releases.FindLatest(constraint)
	if err != nil {
		return nil, fmt.Errorf("determineVersion: error finding the latest release: %w", err)
	}
	if !previousReleaseExists {
		return version.Semver(), nil
	}

	maxPublishedVersion, err := ReleaseVersionFromPublishedVersion(latestRelease.Version)
	if err != nil {
		return nil, fmt.Errorf("determineVersion: error parsing release version: %w", err)
	}

	version, err = version.SetPrereleaseVersion(maxPublishedVersion.PrereleaseVersion() + 1)
	if err != nil {
		return nil, err
	}

	return version.Semver(), nil
}

func releaseType(window string, v *releaseVersion) pivnet.ReleaseType {
	switch window {
	case "rc":
		return "Release Candidate"
	case "beta":
		return "Beta Release"
	case "alpha":
		return "Alpha Release"
	case "ga":
		switch {
		case v.IsMajor():
			return "Major Release"
		case v.IsMinor():
			return "Minor Release"
		default:
			return "Maintenance Release"
		}
	default:
		return "Developer Release"
	}
}

// Usage writes helpful information.
func (p Publish) Usage() jhanda.Usage {
	return jhanda.Usage{
		Description:      "This command prints helpful usage information.",
		ShortDescription: "prints this usage information",
		Flags:            p.Options,
	}
}

const PublishDateFormat = "2006-01-02"

type Date struct {
	time.Time
}

// UnmarshalYAML parses a date in "YYYY-MM-DD" format
func (d *Date) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	now, err := time.ParseInLocation(PublishDateFormat, str, time.UTC)
	if err != nil {
		return err
	}

	d.Time = now
	return nil
}

type Kilnfile struct {
	Slug         string `yaml:"slug"`
	PublishDates struct {
		Beta Date `yaml:"beta"`
		RC   Date `yaml:"rc"`
		GA   Date `yaml:"ga"`
	} `yaml:"publish_dates"`
}

// ReleaseWindow determines the release window based on the current time.
func (kilnfile Kilnfile) ReleaseWindow(currentTime time.Time) (string, error) {
	gaDate := kilnfile.PublishDates.GA
	if currentTime.Equal(gaDate.Time) || currentTime.After(gaDate.Time) {
		return "ga", nil
	}

	rcDate := kilnfile.PublishDates.RC
	if currentTime.Equal(rcDate.Time) || currentTime.After(rcDate.Time) {
		return "rc", nil
	}

	betaDate := kilnfile.PublishDates.Beta
	if currentTime.Equal(betaDate.Time) || currentTime.After(betaDate.Time) {
		return "beta", nil
	}

	return "alpha", nil
}

type releaseSet []pivnet.Release

func (rs releaseSet) Find(version string) (pivnet.Release, error) {
	for _, r := range rs {
		if r.Version == version {
			return r, nil
		}
	}

	return pivnet.Release{}, fmt.Errorf("release with version %s not found", version)
}

func (rs releaseSet) FindLatest(constraint *semver.Constraints) (pivnet.Release, bool, error) {
	var matches []pivnet.Release
	for _, release := range rs {
		v, err := semver.NewVersion(release.Version)
		if err != nil {
			continue
		}

		if constraint.Check(v) {
			matches = append(matches, release)
		}
	}

	if len(matches) == 0 {
		return pivnet.Release{}, false, nil
	}

	sort.Slice(matches, func(i, j int) bool {
		v1 := semver.MustParse(matches[i].Version)
		v2 := semver.MustParse(matches[j].Version)
		return v1.LessThan(v2)
	})

	return matches[len(matches)-1], true, nil
}

type releaseVersion struct {
	semver            semver.Version
	window            string
	prereleaseVersion int
}

func ReleaseVersionFromBuildVersion(baseVersion *semver.Version, window string) (*releaseVersion, error) {
	v2, err := baseVersion.SetPrerelease("")
	if err != nil {
		return nil, fmt.Errorf("ReleaseVersionFromBuildVersion: error clearing prerelease of %q: %w", v2, err)
	}

	rv := &releaseVersion{semver: v2, window: window, prereleaseVersion: 0}

	if window != "ga" {
		rv, err = rv.SetPrereleaseVersion(1)
		if err != nil {
			return nil, fmt.Errorf("ReleaseVersionFromBuildVersion: error setting prerelease of %q to 1: %w", rv, err)
		}
	}
	return rv, nil
}

func ReleaseVersionFromPublishedVersion(versionString string) (*releaseVersion, error) {
	version, err := semver.NewVersion(versionString)
	if err != nil {
		return nil, fmt.Errorf("ReleaseVersionFromPublishedVersion: unable to parse version %q: %w", versionString, err)
	}
	segments := strings.Split(version.Prerelease(), ".")
	if len(segments) != 2 {
		return nil, fmt.Errorf("ReleaseVersionFromPublishedVersion: expected prerelease to have a dot (%q)", version)
	}

	window := segments[0]
	prereleaseVersion, err := strconv.Atoi(segments[len(segments)-1])
	if err != nil {
		return nil, fmt.Errorf("ReleaseVersionFromPublishedVersion: release has malformed prelease version (%s): %w", version, err)
	}

	return &releaseVersion{
		semver:            *version,
		window:            window,
		prereleaseVersion: prereleaseVersion,
	}, nil
}

func (rv releaseVersion) MajorMinorConstraint() (*semver.Constraints, error) {
	return semver.NewConstraint(fmt.Sprintf("~%d.%d.0", rv.semver.Major(), rv.semver.Minor()))
}

func (rv releaseVersion) PrereleaseVersionsConstraint() (*semver.Constraints, error) {
	if rv.IsGA() {
		return nil, fmt.Errorf("can't determine PrereleaseVersionsConstraint for %q, which is GA", rv.semver)
	}
	coreVersion := fmt.Sprintf("%d.%d.%d-%s", rv.semver.Major(), rv.semver.Minor(), rv.semver.Patch(), rv.window)
	constraintStr := fmt.Sprintf(">= %s.0, <= %s.9999", coreVersion, coreVersion)
	return semver.NewConstraint(constraintStr)
}

func (rv releaseVersion) SetPrereleaseVersion(prereleaseVersion int) (*releaseVersion, error) {
	if rv.IsGA() {
		return nil, fmt.Errorf("SetPrereleaseVersion: can't set the prerelease version on a GA version (%q)", rv.String())
	}
	v, err := rv.semver.SetPrerelease(fmt.Sprintf("%s.%d", rv.window, prereleaseVersion))
	if err != nil {
		return nil, fmt.Errorf("SetPrereleaseVersion: couldn't set prerelease: %w", err)
	}
	rv.semver = v
	rv.prereleaseVersion = prereleaseVersion

	return &rv, nil
}

func (rv releaseVersion) IsGA() bool {
	return rv.window == "ga"
}

func (rv releaseVersion) IsMajor() bool {
	return rv.semver.Minor() == 0 && rv.semver.Patch() == 0
}

func (rv releaseVersion) IsMinor() bool {
	return rv.semver.Minor() != 0 && rv.semver.Patch() == 0
}

func (rv releaseVersion) String() string {
	return rv.semver.String()
}

func (rv releaseVersion) Semver() *semver.Version {
	return &rv.semver
}

func (rv releaseVersion) PrereleaseVersion() int {
	return rv.prereleaseVersion
}
