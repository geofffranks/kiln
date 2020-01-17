// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"io"
	"sync"

	"github.com/pivotal-cf/kiln/commands"
	"github.com/pivotal-cf/kiln/release"
)

type ReleaseUploader struct {
	DownloadReleasesStub        func(string, []release.RemoteRelease, int) ([]release.LocalRelease, error)
	downloadReleasesMutex       sync.RWMutex
	downloadReleasesArgsForCall []struct {
		arg1 string
		arg2 []release.RemoteRelease
		arg3 int
	}
	downloadReleasesReturns struct {
		result1 []release.LocalRelease
		result2 error
	}
	downloadReleasesReturnsOnCall map[int]struct {
		result1 []release.LocalRelease
		result2 error
	}
	GetMatchedReleasesStub        func(release.ReleaseRequirementSet) ([]release.RemoteRelease, error)
	getMatchedReleasesMutex       sync.RWMutex
	getMatchedReleasesArgsForCall []struct {
		arg1 release.ReleaseRequirementSet
	}
	getMatchedReleasesReturns struct {
		result1 []release.RemoteRelease
		result2 error
	}
	getMatchedReleasesReturnsOnCall map[int]struct {
		result1 []release.RemoteRelease
		result2 error
	}
	IDStub        func() string
	iDMutex       sync.RWMutex
	iDArgsForCall []struct {
	}
	iDReturns struct {
		result1 string
	}
	iDReturnsOnCall map[int]struct {
		result1 string
	}
	UploadReleaseStub        func(string, string, io.Reader) error
	uploadReleaseMutex       sync.RWMutex
	uploadReleaseArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 io.Reader
	}
	uploadReleaseReturns struct {
		result1 error
	}
	uploadReleaseReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ReleaseUploader) DownloadReleases(arg1 string, arg2 []release.RemoteRelease, arg3 int) ([]release.LocalRelease, error) {
	var arg2Copy []release.RemoteRelease
	if arg2 != nil {
		arg2Copy = make([]release.RemoteRelease, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.downloadReleasesMutex.Lock()
	ret, specificReturn := fake.downloadReleasesReturnsOnCall[len(fake.downloadReleasesArgsForCall)]
	fake.downloadReleasesArgsForCall = append(fake.downloadReleasesArgsForCall, struct {
		arg1 string
		arg2 []release.RemoteRelease
		arg3 int
	}{arg1, arg2Copy, arg3})
	fake.recordInvocation("DownloadReleases", []interface{}{arg1, arg2Copy, arg3})
	fake.downloadReleasesMutex.Unlock()
	if fake.DownloadReleasesStub != nil {
		return fake.DownloadReleasesStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.downloadReleasesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ReleaseUploader) DownloadReleasesCallCount() int {
	fake.downloadReleasesMutex.RLock()
	defer fake.downloadReleasesMutex.RUnlock()
	return len(fake.downloadReleasesArgsForCall)
}

func (fake *ReleaseUploader) DownloadReleasesCalls(stub func(string, []release.RemoteRelease, int) ([]release.LocalRelease, error)) {
	fake.downloadReleasesMutex.Lock()
	defer fake.downloadReleasesMutex.Unlock()
	fake.DownloadReleasesStub = stub
}

func (fake *ReleaseUploader) DownloadReleasesArgsForCall(i int) (string, []release.RemoteRelease, int) {
	fake.downloadReleasesMutex.RLock()
	defer fake.downloadReleasesMutex.RUnlock()
	argsForCall := fake.downloadReleasesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *ReleaseUploader) DownloadReleasesReturns(result1 []release.LocalRelease, result2 error) {
	fake.downloadReleasesMutex.Lock()
	defer fake.downloadReleasesMutex.Unlock()
	fake.DownloadReleasesStub = nil
	fake.downloadReleasesReturns = struct {
		result1 []release.LocalRelease
		result2 error
	}{result1, result2}
}

func (fake *ReleaseUploader) DownloadReleasesReturnsOnCall(i int, result1 []release.LocalRelease, result2 error) {
	fake.downloadReleasesMutex.Lock()
	defer fake.downloadReleasesMutex.Unlock()
	fake.DownloadReleasesStub = nil
	if fake.downloadReleasesReturnsOnCall == nil {
		fake.downloadReleasesReturnsOnCall = make(map[int]struct {
			result1 []release.LocalRelease
			result2 error
		})
	}
	fake.downloadReleasesReturnsOnCall[i] = struct {
		result1 []release.LocalRelease
		result2 error
	}{result1, result2}
}

func (fake *ReleaseUploader) GetMatchedReleases(arg1 release.ReleaseRequirementSet) ([]release.RemoteRelease, error) {
	fake.getMatchedReleasesMutex.Lock()
	ret, specificReturn := fake.getMatchedReleasesReturnsOnCall[len(fake.getMatchedReleasesArgsForCall)]
	fake.getMatchedReleasesArgsForCall = append(fake.getMatchedReleasesArgsForCall, struct {
		arg1 release.ReleaseRequirementSet
	}{arg1})
	fake.recordInvocation("GetMatchedReleases", []interface{}{arg1})
	fake.getMatchedReleasesMutex.Unlock()
	if fake.GetMatchedReleasesStub != nil {
		return fake.GetMatchedReleasesStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getMatchedReleasesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ReleaseUploader) GetMatchedReleasesCallCount() int {
	fake.getMatchedReleasesMutex.RLock()
	defer fake.getMatchedReleasesMutex.RUnlock()
	return len(fake.getMatchedReleasesArgsForCall)
}

func (fake *ReleaseUploader) GetMatchedReleasesCalls(stub func(release.ReleaseRequirementSet) ([]release.RemoteRelease, error)) {
	fake.getMatchedReleasesMutex.Lock()
	defer fake.getMatchedReleasesMutex.Unlock()
	fake.GetMatchedReleasesStub = stub
}

func (fake *ReleaseUploader) GetMatchedReleasesArgsForCall(i int) release.ReleaseRequirementSet {
	fake.getMatchedReleasesMutex.RLock()
	defer fake.getMatchedReleasesMutex.RUnlock()
	argsForCall := fake.getMatchedReleasesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ReleaseUploader) GetMatchedReleasesReturns(result1 []release.RemoteRelease, result2 error) {
	fake.getMatchedReleasesMutex.Lock()
	defer fake.getMatchedReleasesMutex.Unlock()
	fake.GetMatchedReleasesStub = nil
	fake.getMatchedReleasesReturns = struct {
		result1 []release.RemoteRelease
		result2 error
	}{result1, result2}
}

func (fake *ReleaseUploader) GetMatchedReleasesReturnsOnCall(i int, result1 []release.RemoteRelease, result2 error) {
	fake.getMatchedReleasesMutex.Lock()
	defer fake.getMatchedReleasesMutex.Unlock()
	fake.GetMatchedReleasesStub = nil
	if fake.getMatchedReleasesReturnsOnCall == nil {
		fake.getMatchedReleasesReturnsOnCall = make(map[int]struct {
			result1 []release.RemoteRelease
			result2 error
		})
	}
	fake.getMatchedReleasesReturnsOnCall[i] = struct {
		result1 []release.RemoteRelease
		result2 error
	}{result1, result2}
}

func (fake *ReleaseUploader) ID() string {
	fake.iDMutex.Lock()
	ret, specificReturn := fake.iDReturnsOnCall[len(fake.iDArgsForCall)]
	fake.iDArgsForCall = append(fake.iDArgsForCall, struct {
	}{})
	fake.recordInvocation("ID", []interface{}{})
	fake.iDMutex.Unlock()
	if fake.IDStub != nil {
		return fake.IDStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.iDReturns
	return fakeReturns.result1
}

func (fake *ReleaseUploader) IDCallCount() int {
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	return len(fake.iDArgsForCall)
}

func (fake *ReleaseUploader) IDCalls(stub func() string) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = stub
}

func (fake *ReleaseUploader) IDReturns(result1 string) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = nil
	fake.iDReturns = struct {
		result1 string
	}{result1}
}

func (fake *ReleaseUploader) IDReturnsOnCall(i int, result1 string) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = nil
	if fake.iDReturnsOnCall == nil {
		fake.iDReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.iDReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *ReleaseUploader) UploadRelease(arg1 string, arg2 string, arg3 io.Reader) error {
	fake.uploadReleaseMutex.Lock()
	ret, specificReturn := fake.uploadReleaseReturnsOnCall[len(fake.uploadReleaseArgsForCall)]
	fake.uploadReleaseArgsForCall = append(fake.uploadReleaseArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 io.Reader
	}{arg1, arg2, arg3})
	fake.recordInvocation("UploadRelease", []interface{}{arg1, arg2, arg3})
	fake.uploadReleaseMutex.Unlock()
	if fake.UploadReleaseStub != nil {
		return fake.UploadReleaseStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.uploadReleaseReturns
	return fakeReturns.result1
}

func (fake *ReleaseUploader) UploadReleaseCallCount() int {
	fake.uploadReleaseMutex.RLock()
	defer fake.uploadReleaseMutex.RUnlock()
	return len(fake.uploadReleaseArgsForCall)
}

func (fake *ReleaseUploader) UploadReleaseCalls(stub func(string, string, io.Reader) error) {
	fake.uploadReleaseMutex.Lock()
	defer fake.uploadReleaseMutex.Unlock()
	fake.UploadReleaseStub = stub
}

func (fake *ReleaseUploader) UploadReleaseArgsForCall(i int) (string, string, io.Reader) {
	fake.uploadReleaseMutex.RLock()
	defer fake.uploadReleaseMutex.RUnlock()
	argsForCall := fake.uploadReleaseArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *ReleaseUploader) UploadReleaseReturns(result1 error) {
	fake.uploadReleaseMutex.Lock()
	defer fake.uploadReleaseMutex.Unlock()
	fake.UploadReleaseStub = nil
	fake.uploadReleaseReturns = struct {
		result1 error
	}{result1}
}

func (fake *ReleaseUploader) UploadReleaseReturnsOnCall(i int, result1 error) {
	fake.uploadReleaseMutex.Lock()
	defer fake.uploadReleaseMutex.Unlock()
	fake.UploadReleaseStub = nil
	if fake.uploadReleaseReturnsOnCall == nil {
		fake.uploadReleaseReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.uploadReleaseReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ReleaseUploader) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.downloadReleasesMutex.RLock()
	defer fake.downloadReleasesMutex.RUnlock()
	fake.getMatchedReleasesMutex.RLock()
	defer fake.getMatchedReleasesMutex.RUnlock()
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	fake.uploadReleaseMutex.RLock()
	defer fake.uploadReleaseMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ReleaseUploader) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ commands.ReleaseUploader = new(ReleaseUploader)
