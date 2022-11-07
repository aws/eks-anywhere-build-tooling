package builder

import (
	"testing"
)

func TestGetSupportedReleaseBranchesSuccess(t *testing.T) {
	b := BuildOptions{
		ReleaseChannel: "1-24",
	}

	supportedReleaseBranches := GetSupportedReleaseBranches()
	if !SliceContains(supportedReleaseBranches, b.ReleaseChannel) {
		t.Fatalf("GetSupportedReleaseBranches error: supported branches does not contain the release channel"+
			": %s", b.ReleaseChannel)
	}
}

func TestGetSupportedReleaseBranchesFailure(t *testing.T) {
	b := BuildOptions{
		ReleaseChannel: "1-16",
	}

	supportedReleaseBranches := GetSupportedReleaseBranches()
	if SliceContains(supportedReleaseBranches, b.ReleaseChannel) {
		t.Fatalf("GetSupportedReleaseBranches error: supported branches does not contain the release channel"+
			": %s", b.ReleaseChannel)
	}
}
