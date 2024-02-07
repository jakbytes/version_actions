package github

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoReleaseVersionFound_Error(t *testing.T) {
	err := NoReleaseVersionFound{}
	require.Equal(t, "no commits found with a release tag", err.Error())
}

func TestNoPullRequestFoundError_Error(t *testing.T) {
	err := NoPullRequestFoundError{
		Head: "head",
		Base: "base",
	}
	require.Equal(t, "no pull request found for branch head targeting base", err.Error())
}

func TestMultiplePullRequestsFoundError_Error(t *testing.T) {
	err := MultiplePullRequestsFoundError{
		Head: "head",
		Base: "base",
	}
	require.Equal(t, "multiple pull requests found for branch head targeting base", err.Error())
}

func TestBranchNotFound_Error(t *testing.T) {
	err := BranchNotFound{
		Name: "name",
	}
	require.Equal(t, "branch name not found", err.Error())
}

func TestNoPrereleaseVersionFound_Error(t *testing.T) {
	err := NoPrereleaseVersionFound{}
	require.Equal(t, "no commits found with a prerelease tag", err.Error())
}
