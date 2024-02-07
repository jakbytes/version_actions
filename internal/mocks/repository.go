package mocks

import (
	"context"
	"github.com/google/go-github/v58/github"
	"time"
)

type RepositoryService struct {
	Inner          error
	CompareError   error
	GetBranchError func(ctx context.Context, owner string, repo string, branch string, maxRedirects int) error
	GetError       error
	Commits        []*github.RepositoryCommit
	Tags           []*github.RepositoryTag
	Comparison     *github.CommitsComparison
}

func (r *RepositoryService) GetBranch(ctx context.Context, owner string, repo string, branch string, maxRedirects int) (*github.Branch, *github.Response, error) {
	if r.GetBranchError != nil {
		err := r.GetBranchError(ctx, owner, repo, branch, maxRedirects)
		if err != nil {
			return nil, nil, err
		}
	}
	return &github.Branch{
		Name: github.String(branch),
		Commit: &github.RepositoryCommit{
			SHA: github.String("hash"),
		},
	}, nil, nil
}

func (r *RepositoryService) Get(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
	if r.GetError != nil {
		return nil, nil, r.GetError
	}
	return &github.Repository{
		Name:          github.String(repo),
		Owner:         &github.User{Login: github.String(owner)},
		DefaultBranch: github.String("main"),
	}, nil, nil
}

func (r *RepositoryService) CompareCommits(ctx context.Context, owner string, repo string, base string, head string, opts *github.ListOptions) (*github.CommitsComparison, *github.Response, error) {
	if r.Inner != nil {
		return nil, nil, r.Inner
	}
	if r.CompareError != nil {
		return nil, nil, r.CompareError
	}
	if r.Comparison != nil {
		return r.Comparison, nil, nil
	}
	return &github.CommitsComparison{
		Commits: []*github.RepositoryCommit{
			{
				SHA: github.String("hash1-hash1"),
				Commit: &github.Commit{
					Message: github.String("feat: message1"),
					Author: &github.CommitAuthor{
						Login: github.String("login1"),
					},
					Committer: &github.CommitAuthor{
						Login: github.String("user1"),
						Date:  &github.Timestamp{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
					},
				},
			},
			{
				SHA: github.String("hash2-hash2"),
				Commit: &github.Commit{
					Message: github.String("fix: message2"),
					Author: &github.CommitAuthor{
						Login: github.String("login2"),
					},
					Committer: &github.CommitAuthor{
						Login: github.String("user2"),
						Date:  &github.Timestamp{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
					},
				},
			},
			{
				SHA: github.String("hash1-hash1"), // Duplicate hash to test handling of duplicates
				Commit: &github.Commit{
					Message: github.String("feat: message1"),
					Author: &github.CommitAuthor{
						Login: github.String("login1"),
					},
					Committer: &github.CommitAuthor{
						Login: github.String("user1"),
						Date:  &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
					},
				},
			},
			// Add more Commits as necessary
		},
	}, nil, nil
}

func (r *RepositoryService) ListCommits(ctx context.Context, owner string, repo string, opts *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error) {
	if r.Inner != nil {
		return nil, nil, r.Inner
	}
	if r.Commits != nil {
		var commits []*github.RepositoryCommit
		for i := opts.Page; i < opts.Page+opts.PerPage; i++ {
			// If we've reached the end of the Commits, return what we have
			if i >= len(r.Commits) {
				// ensure to return the next page
				return commits, &github.Response{
					NextPage: 0,
				}, nil
			}
			commits = append(commits, r.Commits[i])
		}
		return commits, &github.Response{
			NextPage: opts.Page + opts.PerPage,
		}, nil
	}
	return []*github.RepositoryCommit{
		{
			SHA: github.String("hash1-hash1"),
			Commit: &github.Commit{
				Message: github.String("feat: message1"),
				Committer: &github.CommitAuthor{
					Date: &github.Timestamp{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
					Name: github.String("name"),
				},
				Tree: &github.Tree{
					SHA: github.String("tree1-tree1"),
				},
			},
		},
		{
			SHA:    github.String("hash2-hash2"),
			Commit: &github.Commit{},
		},
		{
			SHA:    github.String("hash3-hash3"),
			Commit: &github.Commit{},
		},
	}, nil, nil
}

func (r *RepositoryService) ListTags(ctx context.Context, owner string, repo string, opts *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error) {
	if r.Inner != nil {
		return nil, nil, r.Inner
	}
	if r.Tags != nil {
		return r.Tags, nil, nil
	}
	return []*github.RepositoryTag{
		{
			Name: github.String("v1.0.0"),
			Commit: &github.Commit{
				SHA: github.String("hash1-hash1"),
			},
		},
		{
			Name: github.String("v1.0.1"),
			Commit: &github.Commit{
				SHA: github.String("hash2-hash2"),
			},
		},
		{
			Name: github.String("v1.0.2-rc.1"),
			Commit: &github.Commit{
				SHA: github.String("hash3-hash3"),
			},
		},
	}, nil, nil
}
