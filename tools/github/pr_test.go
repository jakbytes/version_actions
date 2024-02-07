package github

import (
	"context"
	"fmt"
	"github.com/jakbytes/version_actions/internal/mocks"
	"github.com/jakbytes/version_actions/tools/changelog"
	"strings"
	"testing"

	"github.com/google/go-github/v58/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePullRequest(t *testing.T) {
	client := NewClient(context.Background(), "token", "owner", "name")
	client.PullRequests = &mocks.PullRequestsService{}

	// Test data
	head := "dev"
	base := "main"
	title := "New Feature"
	body := changelog.Markdown{"This is a new feature"}

	pr, err := client.CreatePullRequest(head, base, title, body, false)
	require.Nil(t, err)

	require.Equal(t, title, *pr.Title)
	require.Equal(t, head, *pr.Head.Ref)
	require.Equal(t, base, *pr.Base.Ref)
	require.Equal(t, body[0], *pr.Body)
}

func TestCreatePullRequest_Error(t *testing.T) {
	client := NewClient(context.Background(), "token", "owner", "name")
	client.PullRequests = &mocks.PullRequestsService{Inner: &github.ErrorResponse{}}

	// Test data
	head := "dev"
	base := "main"
	title := "New Feature"
	body := changelog.Markdown{"This is a new feature"}

	pr, err := client.CreatePullRequest(head, base, title, body, false)
	require.NotNil(t, err)
	require.Nil(t, pr)
}

func TestGetPullRequest(t *testing.T) {
	client := NewClient(context.Background(), "token", "owner", "name")
	client.PullRequests = &mocks.PullRequestsService{}

	// Test data
	head := "dev"
	base := "main"

	pr, err := client.GetPullRequest(head, base)
	require.Nil(t, err)

	require.Equal(t, head, *pr.Head.Ref)
	require.Equal(t, base, *pr.Base.Ref)
}

func TestGetPullRequest_Error(t *testing.T) {
	client := NewClient(context.Background(), "token", "owner", "name")
	client.PullRequests = &mocks.PullRequestsService{Inner: &github.ErrorResponse{}}

	// Test data
	head := "dev"
	base := "main"

	pr, err := client.GetPullRequest(head, base)
	require.NotNil(t, err)
	require.Nil(t, pr)
}

func TestGetPullRequest_NotFound(t *testing.T) {
	client := NewClient(context.Background(), "token", "owner", "name")
	client.PullRequests = &mocks.PullRequestsService{
		PullRequests: make([]*github.PullRequest, 0),
	}

	// Test data
	head := "dev"
	base := "main"

	pr, err := client.GetPullRequest(head, base)
	require.NotNil(t, err)
	// passing an empty PR makes some operations with pull requests a little cleaner, this is intentional
	require.NotNil(t, pr)
	require.Nil(t, pr.Body)
	require.Equal(t, NoPullRequestFoundError{Head: head, Base: base}, err)
}

func TestGetPullRequest_MultiplePullRequests(t *testing.T) {
	client := NewClient(context.Background(), "token", "owner", "name")
	client.PullRequests = &mocks.PullRequestsService{
		PullRequests: []*github.PullRequest{
			{},
			{},
		},
	}

	// Test data
	head := "dev"
	base := "main"

	pr, err := client.GetPullRequest(head, base)
	require.NotNil(t, err)
	require.Nil(t, pr)
	require.Equal(t, MultiplePullRequestsFoundError{Head: head, Base: base}, err)
}

func TestEditPullRequest(t *testing.T) {
	client := NewClient(context.Background(), "token", "owner", "name")
	client.PullRequests = &mocks.PullRequestsService{}

	// Test data
	head := "head"
	base := "base"
	title := "New Feature"
	body := changelog.Markdown{"This is a new feature"}

	pr, err := client.EditPullRequest(head, base, title, body)
	require.Nil(t, err)

	require.Equal(t, title, *pr.Title)
	require.Equal(t, head, *pr.Head.Ref)
	require.Equal(t, base, *pr.Base.Ref)
	require.Equal(t, body[0], *pr.Body)
}

func TestEditPullRequest_Error(t *testing.T) {
	client := NewClient(context.Background(), "token", "owner", "name")
	client.PullRequests = &mocks.PullRequestsService{Inner: &github.ErrorResponse{}}

	// Test data
	head := "head"
	base := "base"
	title := "New Feature"
	body := changelog.Markdown{"This is a new feature"}

	pr, err := client.EditPullRequest(head, base, title, body)
	require.NotNil(t, err)
	require.Nil(t, pr)

	client = NewClient(context.Background(), "token", "owner", "name")
	client.PullRequests = &mocks.PullRequestsService{InnerEdit: &github.ErrorResponse{}}

	pr, err = client.EditPullRequest(head, base, title, body)
	require.NotNil(t, err)
	require.Nil(t, pr)
}

func TestSetPullRequest_Create(t *testing.T) {
	client := NewClient(context.Background(), "token", "owner", "name")
	client.PullRequests = &mocks.PullRequestsService{}

	// Test data
	head := "head"
	base := "base"
	title := "New Feature"
	body := changelog.Markdown{"This is a new feature"}

	err := client.SetPullRequest(head, base, title, false, func(_ *string) (changelog.Markdown, error) {
		return body, nil
	})
	require.Nil(t, err)

	pr := client.PullRequests.(*mocks.PullRequestsService).PullRequests[0]

	require.Equal(t, title, *pr.Title)
	require.Equal(t, head, *pr.Head.Ref)
	require.Equal(t, base, *pr.Base.Ref)
	require.Equal(t, body[0], *pr.Body)
}

func TestSetPullRequest_Edit(t *testing.T) {
	client := NewClient(context.Background(), "token", "owner", "name")
	client.PullRequests = &mocks.PullRequestsService{
		PullRequests: []*github.PullRequest{
			{
				Number: github.Int(0),
				Body:   github.String("old body"),
				Head:   &github.PullRequestBranch{Ref: github.String("head")},
				Base:   &github.PullRequestBranch{Ref: github.String("base")},
			},
		},
	}

	// Test data
	head := "head"
	base := "base"
	title := "New Feature"

	err := client.SetPullRequest(head, base, title, false, func(body *string) (changelog.Markdown, error) {
		return strings.Split(*body, "\n"), nil
	})
	require.Nil(t, err)

	pr := client.PullRequests.(*mocks.PullRequestsService).PullRequests[0]

	require.Equal(t, title, *pr.Title)
	require.Equal(t, head, *pr.Head.Ref)
	require.Equal(t, base, *pr.Base.Ref)
	require.Equal(t, "old body", *pr.Body)
}

func TestSetPullRequest_ComposeBodyError(t *testing.T) {
	client := NewClient(context.Background(), "token", "owner", "name")
	client.PullRequests = &mocks.PullRequestsService{}

	// Test data
	head := "head"
	base := "base"
	title := "New Feature"

	err := client.SetPullRequest(head, base, title, false, func(body *string) (changelog.Markdown, error) {
		return changelog.Markdown{""}, assert.AnError
	})
	require.NotNil(t, err)
	require.Equal(t, fmt.Errorf("failed to compose pull request body: %w", assert.AnError), err)
}

func TestSetPullRequest_GetPullRequestError(t *testing.T) {
	client := NewClient(context.Background(), "token", "owner", "name")
	client.PullRequests = &mocks.PullRequestsService{Inner: assert.AnError}

	// Test data
	head := "head"
	base := "base"
	title := "New Feature"

	err := client.SetPullRequest(head, base, title, false, func(body *string) (changelog.Markdown, error) {
		return changelog.Markdown{}, nil
	})
	require.NotNil(t, err)
	require.Equal(t, fmt.Errorf("unable to verify if existing pull request exists: %w", assert.AnError), err)
}

func TestPullRequestBody_String(t *testing.T) {
	body := changelog.Markdown{"line1", "line2"}

	require.Equal(t, "line1\nline2", body.String())
}
