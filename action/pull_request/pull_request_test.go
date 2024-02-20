package pull_request

import (
	"context"
	"github.com/jakbytes/version_actions/internal/mocks"
	"github.com/jakbytes/version_actions/tools/github"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComposePullRequestTitle(t *testing.T) {
	client := github.NewClient(context.Background(), "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{
		// commit with message under 70 characters
		Commits: []*github.RepositoryCommit{
			{
				SHA: github.String("hash1-hash1"),
				Commit: &github.Commit{
					Message: github.String("feat: commit message"),
				},
			},
		},
	}

	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)

	title, err := composeTitle(branch)
	require.Nil(t, err)

	require.Equal(t, "feat: commit message", title)
}

func TestComposePullRequestTitle_Error(t *testing.T) {
	client := github.NewClient(context.Background(), "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{
		Inner: assert.AnError,
	}

	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)

	_, err = composeTitle(branch)
	require.NotNil(t, err)
	require.Equal(t, assert.AnError, err)
}

func TestComposePullRequestTitle_LongMessage(t *testing.T) {
	client := github.NewClient(context.Background(), "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{
		// commit with message over 70 characters
		Commits: []*github.RepositoryCommit{
			{
				SHA: github.String("hash1-hash1"),
				Commit: &github.Commit{
					Message: github.String("feat: commit message with a long message that is over 70 characters, this part gets truncated"),
				},
			},
		},
	}

	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)

	title, err := composeTitle(branch)
	require.Nil(t, err)

	require.Equal(t, "feat: commit message with a long message that is over 70 characters, t...", title)
}

func TestComposePullRequestBody(t *testing.T) {
	client := github.NewClient(context.Background(), "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{}

	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)

	body, err := composeBody(branch, "base", nil)
	require.Nil(t, err)

	expected := []string{
		"### :robot: I have created a pull request *beep* *boop*",
		"",
		"### Notes",
		"",
		"You can add your personal notes here (above the 'Changelog' section). To ensure your notes and the automated " +
			"changelog updates are maintained correctly, keep the 'Changelog' marker in place. If the 'Changelog' marker " +
			"is removed, the automated updates to the changelog will not occur. Personal notes above the 'Changelog' " +
			"will be retained during updates, while content below it will be updated with each new commit.",
		"",
		"## Changelog",
		"### Features",
		"",
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) message1",
		"",
		"### Fixes",
		"",
		"- ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2)) message2",
		"",
		"#",
		"",
		"This Changelog was composed by [version_action](https://github.com/jakbytes/version_action)",
	}

	for i, line := range body {
		assert.Equal(t, expected[i], line)
	}
}

func TestUpdateBody(t *testing.T) {
	existing := github.String(strings.Join([]string{
		"existing body",
		"",
		"## Changelog",
		"### Features",
		"",
		"- message1",
		"  > _Contributed by [](https://github.com/) on 2022-01-01 00:00 UTC_ ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1))",
		"",
	}, "\n"))

	changelog := []string{
		"## Changelog",
		"### Features",
		"",
		"- message1",
		"  > _Contributed by [](https://github.com/) on 2022-01-01 00:00 UTC_ ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1))",
		"",
		"### Fixes",
		"",
		"- message2",
		"  > _Contributed by [](https://github.com/) on 2022-01-02 00:00 UTC_ ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2))",
		"",
		"---",
		"",
		"This Changelog was composed by [version_action](https://github.com/jakbytes/version_action)",
	}

	expected := []string{
		"existing body",
		"",
		"## Changelog",
		"### Features",
		"",
		"- message1",
		"  > _Contributed by [](https://github.com/) on 2022-01-01 00:00 UTC_ ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1))",
		"",
		"### Fixes",
		"",
		"- message2",
		"  > _Contributed by [](https://github.com/) on 2022-01-02 00:00 UTC_ ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2))",
		"",
		"---",
		"",
		"This Changelog was composed by [version_action](https://github.com/jakbytes/version_action)",
	}

	body := updateBody(existing, changelog)

	for i, line := range body {
		assert.Equal(t, expected[i], line)
	}

	existing = github.String(strings.Join([]string{
		"existing body",
		"",
		"",
		"## Changelog",
		"### Features",
		"",
		"- message1",
		"  > _Contributed by [](https://github.com/) on 2022-01-01 00:00 UTC_ ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1))",
		"",
	}, "\n"))

	changelog = []string{
		"## Changelog",
		"### Features",
		"",
		"- message1",
		"  > _Contributed by [](https://github.com/) on 2022-01-01 00:00 UTC_ ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1))",
		"",
		"### Fixes",
		"",
		"- message2",
		"  > _Contributed by [](https://github.com/) on 2022-01-02 00:00 UTC_ ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2))",
		"",
		"---",
		"",
		"This Changelog was composed by [version_action](https://github.com/jakbytes/version_action)",
	}

	expected = []string{
		"existing body",
		"",
		"",
		"## Changelog",
		"### Features",
		"",
		"- message1",
		"  > _Contributed by [](https://github.com/) on 2022-01-01 00:00 UTC_ ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1))",
		"",
		"### Fixes",
		"",
		"- message2",
		"  > _Contributed by [](https://github.com/) on 2022-01-02 00:00 UTC_ ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2))",
		"",
		"---",
		"",
		"This Changelog was composed by [version_action](https://github.com/jakbytes/version_action)",
	}

	body = updateBody(existing, changelog)

	for i, line := range body {
		assert.Equal(t, expected[i], line)
	}

	existing = github.String(strings.Join([]string{
		"",
		"existing body",
		"",
		"",
		"## Changelog",
		"### Features",
		"",
		"- message1",
		"  > _Contributed by [](https://github.com/) on 2022-01-01 00:00 UTC_ ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1))",
		"",
	}, "\n"))

	changelog = []string{
		"## Changelog",
		"### Features",
		"",
		"- message1",
		"  > _Contributed by [](https://github.com/) on 2022-01-01 00:00 UTC_ ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1))",
		"",
		"### Fixes",
		"",
		"- message2",
		"  > _Contributed by [](https://github.com/) on 2022-01-02 00:00 UTC_ ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2))",
		"",
		"---",
		"",
		"This Changelog was composed by [version_action](https://github.com/jakbytes/version_action)",
	}

	expected = []string{
		"",
		"existing body",
		"",
		"",
		"## Changelog",
		"### Features",
		"",
		"- message1",
		"  > _Contributed by [](https://github.com/) on 2022-01-01 00:00 UTC_ ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1))",
		"",
		"### Fixes",
		"",
		"- message2",
		"  > _Contributed by [](https://github.com/) on 2022-01-02 00:00 UTC_ ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2))",
		"",
		"---",
		"",
		"This Changelog was composed by [version_action](https://github.com/jakbytes/version_action)",
	}

	body = updateBody(existing, changelog)

	for i, line := range body {
		assert.Equal(t, expected[i], line)
	}
}

func TestComposePullRequestBody_ExistingBody(t *testing.T) {
	client := github.NewClient(context.Background(), "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{}

	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)

	body, err := composeBody(branch, "base", github.String(strings.Join([]string{
		"existing body",
		"",
		"## Changelog",
		"### Features",
		"",
		"- message1",
		"  > _Contributed by [](https://github.com/) on 2022-01-01 00:00 UTC_ ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1))",
		"",
	}, "\n")))
	require.Nil(t, err)

	expected := []string{
		"existing body",
		"",
		"## Changelog",
		"### Features",
		"",
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) message1",
		"",
		"### Fixes",
		"",
		"- ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2)) message2",
		"",
		"---",
		"",
		"This Changelog was composed by [version_action](https://github.com/jakbytes/version_action)",
	}

	for i, line := range body {
		assert.Equal(t, expected[i], line)
	}
}

func TestComposePullRequestBody_ExistingBodyNoChangelog(t *testing.T) {
	client := github.NewClient(context.Background(), "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{}

	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)

	body, err := composeBody(branch, "base", github.String("existing body"))
	require.Nil(t, err)

	expected := []string{
		"existing body",
	}

	for i, line := range body {
		assert.Equal(t, expected[i], line)
	}
}

func TestComposePullRequestBody_Error(t *testing.T) {
	client := github.NewClient(context.Background(), "token", "owner", "name")
	client.Repositories = &mocks.RepositoryService{
		Inner: assert.AnError,
	}

	branch, err := client.Repository().Branch("branch")
	require.Nil(t, err)

	_, err = composeBody(branch, "base", nil)
	require.NotNil(t, err)
	require.Equal(t, assert.AnError, err)
}

func TestGetArgsValid(t *testing.T) {
	// Set up
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base"}

	// Execute
	result := getArgs()

	// Assert
	expected := Args{
		Action: "action",
		Token:  "token",
		Owner:  "owner",
		Name:   "name",
		Head:   "head",
		Base:   "base",
	}
	assert.Equal(t, expected, result, "The two structs should be equal")
}

func TestGetArgsInvalid(t *testing.T) {
	// Set up
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	os.Args = []string{"program", "action"}

	// Execute and Assert
	assert.Panics(t, func() { getArgs() }, "The function should panic for invalid args")
}

/*
func TestSetPullRequest(t *testing.T) {
	prs := &mocks.PullRequestsService{}
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Repositories: &mocks.RepositoryService{},
			PullRequests: prs,
			RepositoryMetadata: github.RepositoryMetadata{
				Owner: owner,
				Name:  name,
			},
		}
	}

	// Set up command line arguments
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base"}

	// Run the function
	assert.NotPanics(t, func() {
		Execute()
	}, "setPullRequest should not panic on successful run")

	require.Equal(t, 1, len(prs.PullRequests))
	pr := prs.PullRequests[0]
	require.Equal(t, "feat: message1", *pr.Title)
	require.Equal(t, "head", *pr.Head.Ref)
	require.Equal(t, "base", *pr.Base.Ref)
	body := *pr.Body

	expected := []string{
		":robot: I have created a pull request *beep* *boop*",
		"",
		"---",
		"",
		"### Notes",
		"",
		"---",
		"",
		"You can add your personal notes here (above the 'Changelog' section). To ensure your notes and the automated " +
			"changelog updates are maintained correctly, keep the 'Changelog' marker in place. If the 'Changelog' marker " +
			"is removed, the automated updates to the changelog will not occur. Personal notes above the 'Changelog' " +
			"will be retained during updates, while content below it will be updated with each new commit.",
		"",
		"## Changelog",
		"### Features",
		"",
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) message1",
		"",
		"### Fixes",
		"",
		"- ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2)) message2",
		"",
		"---",
		"",
		"This Changelog was composed by [version_action](https://github.com/jakbytes/version_action)",
	}

	for i, line := range strings.Split(body, "\n") {
		assert.Equal(t, expected[i], line)
	}
}

func TestSetPullRequest_ExistingPullRequest(t *testing.T) {
	prs := &mocks.PullRequestsService{
		// Return a pull request
		PullRequests: []*github.PullRequest{
			{
				Number: github.Int(0),
				Title:  github.String("feat: message1"),
				Body:   github.String("existing body\n\n## Changelog"),
				Head:   &github.PullRequestBranch{Ref: github.String("head")},
				Base:   &github.PullRequestBranch{Ref: github.String("base")},
			},
		},
	}
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Repositories: &mocks.RepositoryService{},
			PullRequests: prs,
			RepositoryMetadata: github.RepositoryMetadata{
				Owner: owner,
				Name:  name,
			},
		}
	}

	// Set up command line arguments
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base"}

	// Run the function
	err := setPullRequest()
	require.Nil(t, err)

	require.True(t, len(prs.PullRequests) == 1, "There should be one pull request")
	pr := prs.PullRequests[0]
	require.Equal(t, "feat: message1", *pr.Title)
	require.Equal(t, "head", *pr.Head.Ref)
	require.Equal(t, "base", *pr.Base.Ref)
	body := *pr.Body

	expected := []string{
		"existing body",
		"",
		"## Changelog",
		"### Features",
		"",
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) message1",
		"",
		"### Fixes",
		"",
		"- ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2)) message2",
		"",
		"---",
		"",
		"This Changelog was composed by [version_action](https://github.com/jakbytes/version_action)",
	}

	for i, line := range strings.Split(body, "\n") {
		assert.Equal(t, expected[i], line)
	}
}


*/
/*
func TestSetPullRequest_FailedToComposePullRequestTitleError(t *testing.T) {
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Repositories: &mocks.RepositoryService{
				Inner: assert.AnError,
			},
			PullRequests: &mocks.PullRequestsService{},
		}
	}

	// Set up command line arguments
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base"}

	// Run the function
	assert.Panics(t, func() {
		Execute()
	}, "setPullRequest should not panic on error")

	err := setPullRequest()
	require.NotNil(t, err)
	require.Equal(t, fmt.Errorf("failed to compose pull request title: %w", assert.AnError), err)
}



func TestSetPullRequest_FailedToGetBranchError(t *testing.T) {
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Repositories: &mocks.RepositoryService{
				GetBranchError: func(ctx context.Context, owner string, repo string, branch string, maxRedirects int) error {
					return assert.AnError
				},
			},
			PullRequests: &mocks.PullRequestsService{},
		}
	}

	// Set up command line arguments
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base"}

	// Run the function
	err := setPullRequest()
	require.NotNil(t, err)
	require.Equal(t, fmt.Errorf("failed to get branch: %w", assert.AnError), err)
}

func TestSetPullRequest_UnableToVerifyExistingPullRequestExistsError(t *testing.T) {
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Repositories: &mocks.RepositoryService{},
			PullRequests: &mocks.PullRequestsService{
				Inner: assert.AnError,
			},
		}
	}

	// Set up command line arguments
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base"}

	err := setPullRequest()
	require.NotNil(t, err)
	require.Equal(t, fmt.Errorf("unable to verify if existing pull request exists: %w", assert.AnError), err)
}

func TestSetPullRequest_FailedToComposePullRequestBodyError(t *testing.T) {
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Repositories: &mocks.RepositoryService{
				CompareError: assert.AnError,
			},
			PullRequests: &mocks.PullRequestsService{},
		}
	}

	// Set up command line arguments
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base"}

	err := setPullRequest()
	require.NotNil(t, err)
	require.Equal(t, fmt.Errorf("failed to compose pull request body: %w", assert.AnError), err)
}

func TestSetPullRequest_FailedToEditPullRequestError(t *testing.T) {
	prs := &mocks.PullRequestsService{
		// Return a pull request
		PullRequests: []*github.PullRequest{
			{
				Number: github.Int(0),
				Title:  github.String("feat: message1"),
				Body:   github.String("existing body\n## Changelog"),
				Head:   &github.PullRequestBranch{Ref: github.String("head")},
				Base:   &github.PullRequestBranch{Ref: github.String("base")},
			},
		},
		InnerEdit: assert.AnError,
	}
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Repositories: &mocks.RepositoryService{},
			PullRequests: prs,
		}
	}

	// Set up command line arguments
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base"}

	err := setPullRequest()
	require.NotNil(t, err)
	require.Equal(t, assert.AnError, err)
}
*/
