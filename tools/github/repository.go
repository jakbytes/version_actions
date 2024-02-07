package github

import (
	"context"
	"github.com/google/go-github/v58/github"
	"strings"
)

// RepositoriesService is an interface abstracting operations supported by go-github's github.RepositoriesService
type RepositoriesService interface {
	CompareCommits(ctx context.Context, owner string, repo string, base string, head string, opts *github.ListOptions) (*github.CommitsComparison, *github.Response, error)
	ListCommits(ctx context.Context, owner string, repo string, opt *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error)
	ListTags(ctx context.Context, owner string, repo string, opt *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error)
	GetBranch(ctx context.Context, owner string, repo string, branch string, maxRedirects int) (*github.Branch, *github.Response, error)
	Get(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error)
}

// Repository is a struct that contains the RepositoriesService, context, token, owner, and name. It is used to
// interact with the GitHub API in the context of a specific repository.
type Repository struct {
	GitService
	RepositoriesService
	RepositoryMetadata RepositoryMetadata
	repository         *github.Repository
	branches           map[string]*Branch
	Ctx                context.Context
	tags               []*github.RepositoryTag
	versions           *RepositoryVersions
}

func (r *Repository) getBranch(name string) (*github.Branch, error) {
	branch, _, err := r.GetBranch(r.Ctx, r.RepositoryMetadata.Owner, r.RepositoryMetadata.Name, name, 2)
	if err != nil && strings.Contains(err.Error(), "404") {
		return nil, BranchNotFound{Name: name}
	}
	return branch, err
}

// Branch returns a Branch struct that contains the GitHub client, context, token, owner, name, and branch. After
// validating that the branch exists, it will return a Branch struct that can be used to interact with the GitHub API.
func (r *Repository) Branch(name string) (b *Branch, err error) {
	var ok bool
	if b, ok = r.branches[name]; !ok {
		var branch *github.Branch
		branch, err = r.getBranch(name)
		if err != nil {
			return nil, err
		}

		b = &Branch{
			GitService:          r.GitService,
			RepositoriesService: r.RepositoriesService,
			RepositoryMetadata:  r.RepositoryMetadata,
			Branch:              branch,
			Ctx:                 r.Ctx,
			Name:                name,
		}
		r.branches[name] = b
	}
	return b, nil
}

func (r *Repository) DefaultBranch() (branch *Branch, err error) {
	if r.repository == nil {
		r.repository, _, err = r.Get(r.Ctx, r.RepositoryMetadata.Owner, r.RepositoryMetadata.Name)
		if err != nil {
			return nil, err
		}
	}
	return r.Branch(r.repository.GetDefaultBranch())
}

// CreateBranch creates a branch with the given name and SHA. If the SHA is nil, the default branch's SHA will be used.
func (r *Repository) CreateBranch(name string, sha *string) (branch *Branch, err error) {
	newBranch := &github.Reference{Ref: github.String("refs/heads/" + name), Object: &github.GitObject{SHA: sha}}
	_, _, err = r.CreateRef(r.Ctx, r.RepositoryMetadata.Owner, r.RepositoryMetadata.Name, newBranch)
	if err != nil {
		return nil, err
	}
	return r.Branch(name)
}
