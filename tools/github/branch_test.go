package github

import (
	"context"
	"fmt"
	"github.com/jakbytes/version_actions/internal/mocks"
	"testing"

	"github.com/google/go-github/v58/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDistinctCommits(t *testing.T) {
	ctx := context.Background()
	client := NewClient(ctx, "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{}
	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)

	uniqueCommits, err := branch.GetDistinctCommits("base-branch")

	assert.Nil(t, err)
	assert.NotNil(t, uniqueCommits)
	assert.Len(t, uniqueCommits, 2) // Expecting 2 unique commits
}

func TestGetDistinctCommits_Error(t *testing.T) {
	ctx := context.Background()
	client := NewClient(ctx, "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{}
	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)
	branch.RepositoriesService = &mocks.RepositoryService{Inner: assert.AnError}

	uniqueCommits, err := branch.GetDistinctCommits("base-branch")

	assert.NotNil(t, err)
	assert.Nil(t, uniqueCommits)
}

func TestGetCommitsSinceCommit(t *testing.T) {
	ctx := context.Background()
	client := NewClient(ctx, "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{}
	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)
	branch.RepositoriesService = &mocks.RepositoryService{}

	commits, err := branch.GetCommitsSinceCommit(github.String("hash3-hash3"))

	assert.Nil(t, err)
	assert.NotNil(t, commits)
	assert.Len(t, commits, 2) // Expecting 2 commits
}

func TestGetCommitsSinceCommit_Error(t *testing.T) {
	ctx := context.Background()
	client := NewClient(ctx, "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{}
	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)
	branch.RepositoriesService = &mocks.RepositoryService{Inner: assert.AnError}

	commits, err := branch.GetCommitsSinceCommit(github.String("hash3-hash3"))

	assert.NotNil(t, err)
	assert.Nil(t, commits)
}

func TestGetCommitsSinceCommit_Pagination(t *testing.T) {
	// Generate sequence of 15 commits
	commits := make([]*github.RepositoryCommit, 15)
	for i := 0; i < 15; i++ {
		commits[i] = &github.RepositoryCommit{
			SHA:    github.String(fmt.Sprintf("hash%d", i)),
			Commit: &github.Commit{},
		}
	}

	ctx := context.Background()
	client := NewClient(ctx, "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{}
	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)
	branch.RepositoriesService = &mocks.RepositoryService{Commits: commits}

	// nil commit should return all commits
	allCommits, err := branch.GetCommitsSinceCommit(nil)
	assert.Nil(t, err)
	assert.NotNil(t, allCommits)
	assert.Len(t, allCommits, 15)
}

func TestGetLastCommitMessage(t *testing.T) {
	ctx := context.Background()
	client := NewClient(ctx, "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{}
	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)
	branch.RepositoriesService = &mocks.RepositoryService{}

	message, err := branch.GetLastCommitMessage()

	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "feat: message1", message)
}

func TestGetLastCommitMessage_Error(t *testing.T) {
	ctx := context.Background()
	client := NewClient(ctx, "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{}
	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)
	branch.RepositoriesService = &mocks.RepositoryService{Inner: assert.AnError}

	message, err := branch.GetLastCommitMessage()

	assert.NotNil(t, err)
	assert.Empty(t, message)
}

func TestGetLastCommitMessage_NoCommits(t *testing.T) {
	ctx := context.Background()
	client := NewClient(ctx, "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{}
	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)
	branch.RepositoriesService = &mocks.RepositoryService{Commits: []*github.RepositoryCommit{}}

	message, err := branch.GetLastCommitMessage()

	assert.NotNil(t, err)
	assert.Equal(t, "no commits found in branch", err.Error())
	assert.Empty(t, message)
}

func TestReset(t *testing.T) {
	ctx := context.Background()
	client := NewClient(ctx, "token", "owner", "name")
	client.Git = &mocks.GitService{}
	client.Repositories = &mocks.RepositoryService{}
	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)

	err = branch.Reset(github.String("hash"))

	assert.Nil(t, err)
}

func TestReset_Error(t *testing.T) {
	ctx := context.Background()
	client := NewClient(ctx, "token", "owner", "name")
	client.Git = &mocks.GitService{UpdateRefError: assert.AnError}
	client.Repositories = &mocks.RepositoryService{}
	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)

	err = branch.Reset(github.String("hash"))

	assert.NotNil(t, err)
	assert.Error(t, assert.AnError, err)
}
