package mocks

import (
	"context"
	"github.com/google/go-github/v58/github"
	"strings"
)

type PullRequestsService struct {
	Inner        error
	InnerEdit    error
	PullRequests []*github.PullRequest
}

var prs = []*github.PullRequest{
	{
		Number: github.Int(0),
		Title:  github.String(""),
		Head: &github.PullRequestBranch{
			Label: github.String("owner:head"),
			Ref:   github.String("head"),
			SHA:   github.String("sha1"),
			Repo:  nil,
			User:  nil,
		},
		Base: &github.PullRequestBranch{
			Label: github.String("owner:base"),
			Ref:   github.String("base"),
			SHA:   github.String("sha2"),
			Repo:  nil,
			User:  nil,
		},
		Body:  github.String(""),
		Draft: github.Bool(true),
	},
}

func (m *PullRequestsService) Create(ctx context.Context, owner string, repo string, pr *github.NewPullRequest) (*github.PullRequest, *github.Response, error) {
	if m.Inner != nil {
		return nil, nil, m.Inner
	}

	// Mock response - you should tailor this to match what you expect
	mockPR := &github.PullRequest{
		Number: github.Int(len(m.PullRequests)),
		Title:  pr.Title,
		Head: &github.PullRequestBranch{
			Label: github.String(owner + ":" + *pr.Head),
			Ref:   pr.Head,
			SHA:   github.String("sha1"),
			Repo:  nil,
			User:  nil,
		},
		Base: &github.PullRequestBranch{
			Label: github.String(owner + ":" + *pr.Base),
			Ref:   pr.Base,
			SHA:   github.String("sha2"),
			Repo:  nil,
			User:  nil,
		},
		Body:  pr.Body,
		Draft: pr.Draft,
	}

	m.PullRequests = append(m.PullRequests, mockPR)

	// Mocking a successful response
	return mockPR, &github.Response{}, nil
}

func (m *PullRequestsService) List(ctx context.Context, owner string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
	if m.Inner != nil {
		return nil, nil, m.Inner
	}

	if m.PullRequests != nil {
		return m.PullRequests, &github.Response{}, nil
	}
	// Mock response - you should tailor this to match what you expect
	mockPR := &github.PullRequest{
		Title:  nil,
		Number: github.Int(0),
		Head: &github.PullRequestBranch{
			Label: github.String(opts.Head),
			Ref:   github.String(strings.Split(opts.Head, ":")[1]),
			SHA:   github.String("sha1"),
			Repo:  nil,
			User:  nil,
		},
		Base: &github.PullRequestBranch{
			Label: github.String(owner + ":" + opts.Base),
			Ref:   github.String(opts.Base),
			SHA:   github.String("sha2"),
			Repo:  nil,
			User:  nil,
		},
		Body:  nil,
		Draft: nil,
	}
	// Mocking a successful response
	return []*github.PullRequest{mockPR}, &github.Response{}, nil
}

func (m *PullRequestsService) Edit(ctx context.Context, owner string, repo string, number int, pr *github.PullRequest) (*github.PullRequest, *github.Response, error) {
	if m.InnerEdit != nil {
		return nil, nil, m.InnerEdit
	}

	if m.PullRequests != nil && len(m.PullRequests) > number {
		m.PullRequests[number].Title = pr.Title
		m.PullRequests[number].Body = pr.Body
		return m.PullRequests[number], &github.Response{}, nil
	}

	// Mock response - you should tailor this to match what you expect
	mockPR := prs[number]
	mockPR.Title = pr.Title
	mockPR.Body = pr.Body

	// Mocking a successful response
	return mockPR, &github.Response{}, nil
}
