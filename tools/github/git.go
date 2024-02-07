package github

import (
	"context"
	"github.com/google/go-github/v58/github"
)

type GitService interface {
	CreateRef(ctx context.Context, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error)
	UpdateRef(ctx context.Context, owner string, repo string, ref *github.Reference, force bool) (*github.Reference, *github.Response, error)
	CreateBlob(ctx context.Context, owner string, repo string, blob *github.Blob) (*github.Blob, *github.Response, error)
	CreateTree(ctx context.Context, owner string, repo string, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error)
	CreateCommit(ctx context.Context, owner string, repo string, commit *Commit, opts *github.CreateCommitOptions) (*Commit, *github.Response, error)
}
