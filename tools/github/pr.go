package github

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v58/github"
	"version_actions/tools/changelog"
)

// PullRequestsService is an interface abstracting operations supported by go-github's github.PullRequestsService
type PullRequestsService interface {
	Create(ctx context.Context, owner string, repo string, newPR *github.NewPullRequest) (*github.PullRequest, *github.Response, error)
	List(ctx context.Context, owner string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
	Edit(ctx context.Context, owner string, repo string, number int, pull *github.PullRequest) (*github.PullRequest, *github.Response, error)
}

func (c *Client) CreatePullRequest(head, base string, title string, body changelog.Markdown, draft bool) (*github.PullRequest, error) {
	newPR := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(head),
		Base:  github.String(base),
		Body:  github.String(body.String()),
		Draft: github.Bool(draft),
	}

	pr, _, err := c.PullRequests.Create(c.Ctx, c.Owner, c.Name, newPR)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (c *Client) GetPullRequest(head string, base string) (*github.PullRequest, error) {
	opts := &github.PullRequestListOptions{
		State: "open",
		Head:  c.Owner + ":" + head,
		Base:  base,
	}
	prs, _, err := c.PullRequests.List(c.Ctx, c.Owner, c.Name, opts)
	if err != nil {
		return nil, err
	}
	if len(prs) == 0 {
		return &github.PullRequest{}, NoPullRequestFoundError{Head: head, Base: base}
	} else if len(prs) > 1 {
		return nil, MultiplePullRequestsFoundError{Head: head, Base: base}
	}
	return prs[0], nil
}

func (c *Client) EditPullRequest(head, base, title string, body changelog.Markdown) (*github.PullRequest, error) {
	pr, err := c.GetPullRequest(head, base)
	if err != nil {
		return nil, err
	}

	// Update the pull request
	update := &github.PullRequest{
		Title: github.String(title),
		Body:  github.String(body.String()),
	}
	updatedPR, _, err := c.PullRequests.Edit(c.Ctx, c.Owner, c.Name, *pr.Number, update)
	if err != nil {
		return nil, err
	}

	return updatedPR, nil
}

func (c *Client) SetPullRequest(head, base, title string, draft bool, composeBody func(body *string) (changelog.Markdown, error)) error {
	pr, err := c.GetPullRequest(head, base)
	if err != nil && !errors.Is(err, NoPullRequestFoundError{Head: head, Base: base}) {
		return fmt.Errorf("unable to verify if existing pull request exists: %w", err)
	}

	body, err := composeBody(pr.Body)
	if err != nil {
		return fmt.Errorf("failed to compose pull request body: %w", err)
	}

	if pr.Body == nil {
		_, err = c.CreatePullRequest(head, base, title, body, draft)
	} else {
		_, err = c.EditPullRequest(head, base, title, body)
	}
	return err
}
