package github

import (
	"context"
	"github.com/google/go-github/v58/github"
	"github.com/jakbytes/version_actions/internal/logger"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type RepositoryCommit = github.RepositoryCommit
type Commit = github.Commit
type PullRequest = github.PullRequest
type PullRequestBranch = github.PullRequestBranch
type RepositoryTag = github.RepositoryTag
type CommitsComparison = github.CommitsComparison
type CommitAuthor = github.CommitAuthor
type Timestamp = github.Timestamp
type Tree = github.Tree

// Client is a struct that contains the go-github client and the repository metadata to interact with the GitHub API.
type Client struct {
	Ctx          context.Context
	PullRequests PullRequestsService
	Repositories RepositoriesService
	Git          GitService
	RepositoryMetadata
}

type RepositoryMetadata struct {
	Owner string
	Name  string
}

func NewClient(ctx context.Context, token, owner, name string) *Client {
	log.Logger = logger.Base()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &Client{
		Ctx:          ctx,
		PullRequests: client.PullRequests,
		Repositories: client.Repositories,
		Git:          client.Git,
		RepositoryMetadata: RepositoryMetadata{
			Owner: owner,
			Name:  name,
		},
	}
}

func (c *Client) Repository() *Repository {
	return &Repository{
		GitService:          c.Git,
		RepositoriesService: c.Repositories,
		RepositoryMetadata:  c.RepositoryMetadata,
		Ctx:                 c.Ctx,
		branches:            make(map[string]*Branch),
	}
}

func String(s string) *string {
	return &s
}

func Bool(b bool) *bool {
	return &b
}

func Int(i int) *int {
	return &i
}
