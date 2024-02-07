package mocks

import (
	"context"
	"github.com/google/go-github/v58/github"
)

type GitService struct {
	CreateRefError error
	UpdateRefError error
}

func (g GitService) CreateBlob(ctx context.Context, owner string, repo string, blob *github.Blob) (*github.Blob, *github.Response, error) {
	return &github.Blob{
		SHA: github.String("hash4-hash4"),
	}, nil, nil
}

func (g GitService) CreateTree(ctx context.Context, owner string, repo string, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error) {
	return &github.Tree{
		SHA: github.String("hash4-hash4"),
	}, nil, nil
}

func (g GitService) CreateCommit(ctx context.Context, owner string, repo string, commit *github.Commit, opts *github.CreateCommitOptions) (*github.Commit, *github.Response, error) {
	return &github.Commit{
		SHA: github.String("hash4-hash4"),
	}, nil, nil
}

func (g GitService) UpdateRef(ctx context.Context, owner string, repo string, ref *github.Reference, force bool) (*github.Reference, *github.Response, error) {

	if g.UpdateRefError != nil {
		return nil, nil, g.UpdateRefError
	}
	return nil, nil, nil
}

func (g GitService) CreateRef(ctx context.Context, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error) {
	if g.CreateRefError != nil {
		return nil, nil, g.CreateRefError
	}
	return nil, nil, nil
}
