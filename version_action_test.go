package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
	"version_actions/action/pull_request"
	"version_actions/action/version"
	"version_actions/internal/mocks"
	"version_actions/tools/changelog"
	"version_actions/tools/github"
)

var repositories = &mocks.RepositoryService{}
var git = &mocks.GitService{}
var prs = &mocks.PullRequestsService{}

const featureBranch = "jak/feature/branch"
const devBranch = "development"
const stagingBranch = "staging"
const mainBranch = "main"
const devPrereleaseIdentifier = "drc"
const stagingPrereleaseIdentifier = "src"
const repositoryName = "name"
const repositoryOwner = "owner"
const devReleaseBranch = "release--branch--development"
const devPromoteBranch = "release--branch--staging"
const stagingPromoteBranch = "release--branch--main"

func newClient(ctx context.Context, token string, owner string, name string) *github.Client {
	return &github.Client{
		Repositories: repositories,
		Git:          git,
		PullRequests: prs,
		RepositoryMetadata: github.RepositoryMetadata{
			Name:  name,
			Owner: owner,
		},
	}
}

func TestAction(t *testing.T) {
	changelog.Path = "test_CHANGELOG.md"
	defer os.Remove(changelog.Path)
	// starting from repository with no tags and two branches main, development
	version.NewClient = newClient
	pull_request.NewClient = newClient

	t.Run("Testing feature branch generated PR", testFeaturePR)
	t.Run("Testing version release candidate", testVersionRC)
	t.Run("Testing prerelease promotion", testPrereleasePromotion)
	t.Run("Testing release promotion", testPromotion)

}

func testFeaturePR(t *testing.T) {
	// the first commit made to the feature branch
	repositories.Commits = []*github.RepositoryCommit{
		{
			SHA: github.String("hash1-hash1"),
			Commit: &github.Commit{
				Message: github.String("feat: init"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
		},
	}
	// the difference in commits between the feature branch and the development branch
	repositories.Comparison = &github.CommitsComparison{
		Commits: []*github.RepositoryCommit{
			{
				SHA: github.String("hash1-hash1"),
				Commit: &github.Commit{
					Message: github.String("feat: init"),
					Author: &github.CommitAuthor{
						Login: github.String("login1"),
					},
					Committer: &github.CommitAuthor{
						Login: github.String("user1"),
						Date:  &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
					},
				},
			},
		},
	}

	prs.PullRequests = []*github.PullRequest{}

	// starting from repository with no tags and three branches main, staging, development

	os.Args = []string{"program", "action", "token", repositoryOwner, repositoryName, featureBranch, devBranch}
	assert.NotPanics(t, func() {
		pull_request.Execute()
	}, "pull request generation should not panic")

	require.Equal(t, 1, len(prs.PullRequests))
	pr := prs.PullRequests[0]
	require.Equal(t, "feat: init", *pr.Title)
	require.Equal(t, featureBranch, *pr.Head.Ref)
	require.Equal(t, devBranch, *pr.Base.Ref)
	require.Equal(t, true, *pr.Draft)

	expected := []string{
		":robot: I have created a pull request *beep* *boop*",
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
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) init",
		"",
		"---",
		"",
		"This Changelog was composed by [version-action](https://github.com/jakbytes/version-action)",
	}

	for i, line := range strings.Split(*pr.Body, "\n") {
		assert.Equal(t, expected[i], line)
	}

	t.Run("Testing updating PR", testUpdateFeaturePR)
}

func testUpdateFeaturePR(t *testing.T) {
	require.Equal(t, 1, len(prs.PullRequests))

	repositories.Commits = []*github.RepositoryCommit{
		{
			SHA: github.String("hash2-hash2"),
			Commit: &github.Commit{
				Message: github.String("fix: fix related to the feature"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				},
				Tree: &github.Tree{
					SHA: github.String("tree1-tree1"),
				},
			},
		},
		{
			SHA: github.String("hash1-hash1"),
			Commit: &github.Commit{
				Message: github.String("feat: init"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
		},
	}

	repositories.Comparison.Commits = repositories.Commits

	require.Equal(t, 1, len(prs.PullRequests))

	assert.NotPanics(t, func() {
		pull_request.Execute()
	}, "pull request generation should not panic")

	require.Equal(t, 1, len(prs.PullRequests))
	pr := prs.PullRequests[0]
	require.Equal(t, "feat: init", *pr.Title)
	require.Equal(t, featureBranch, *pr.Head.Ref)
	require.Equal(t, devBranch, *pr.Base.Ref)
	require.Equal(t, true, *pr.Draft)

	expected := []string{
		":robot: I have created a pull request *beep* *boop*",
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
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) init",
		"",
		"### Fixes",
		"",
		"- ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2)) fix related to the feature",
		"",
		"---",
		"",
		"This Changelog was composed by [version-action](https://github.com/jakbytes/version-action)",
	}

	for i, line := range strings.Split(*pr.Body, "\n") {
		assert.Equal(t, expected[i], line)
	}

	require.True(t, len(prs.PullRequests) == 1, "There should be one pull request")
}

func testVersionRC(t *testing.T) {
	prs.PullRequests = []*github.PullRequest{}    // pull request was merged into development
	repositories.Tags = []*github.RepositoryTag{} // no tags in the repository
	os.Args = []string{"program", "release", "token", repositoryOwner, repositoryName, devBranch, devBranch, devPrereleaseIdentifier, mainBranch, "none"}

	assert.NotPanics(t, func() {
		version.Execute()
	}, "version generation should not panic")

	require.Equal(t, 1, len(prs.PullRequests))
	pr := prs.PullRequests[0]
	require.Equal(t, "release(development): v0.0.0-drc.0", *pr.Title)
	require.Equal(t, devReleaseBranch, *pr.Head.Ref)
	require.Equal(t, devBranch, *pr.Base.Ref)
	require.Equal(t, false, *pr.Draft)

	expected := []string{
		":robot: I have created a release candidate *beep* *boop*",
		"",
		"---",
		"",
		"## [v0.0.0-drc.0] Initial Version _2022-01-01 00:00 UTC_",
		"### Features",
		"",
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) init",
		"",
		"### Fixes",
		"",
		"- ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2)) fix related to the feature",
		"",
		"---",
		"",
		"This release was composed by [version_actions](https://github.com/jakbytes/version_actions)",
	}

	for i, line := range strings.Split(*pr.Body, "\n") {
		if strings.HasPrefix(line, "## [") {
			prefix := strings.Split(line, " (")[0]
			assert.True(t, strings.HasPrefix(expected[i], prefix))
		} else {
			assert.Equal(t, expected[i], line)
		}
	}

	t.Run("Testing updating version release candidate", testUpdateVersionRC)
}

func testUpdateVersionRC(t *testing.T) {
	repositories.Tags = []*github.RepositoryTag{} // no tags in the repository
	repositories.Commits = []*github.RepositoryCommit{
		{
			SHA: github.String("hash3-hash3"),
			Commit: &github.Commit{
				Message: github.String("fix: another fix related to the feature"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
				},
				Tree: &github.Tree{
					SHA: github.String("tree1-tree1"),
				},
			},
		},
		{
			SHA: github.String("hash2-hash2"),
			Commit: &github.Commit{
				Message: github.String("fix: fix related to the feature"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				},
				Tree: &github.Tree{
					SHA: github.String("tree1-tree1"),
				},
			},
		},
		{
			SHA: github.String("hash1-hash1"),
			Commit: &github.Commit{
				Message: github.String("feat: init"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
		},
	}

	repositories.Comparison.Commits = repositories.Commits

	os.Args = []string{"program", "release", "token", repositoryOwner, repositoryName, devBranch, devBranch, devPrereleaseIdentifier, mainBranch, "none"}

	assert.NotPanics(t, func() {
		version.Execute()
	}, "version generation should not panic")

	require.Equal(t, 1, len(prs.PullRequests))
	pr := prs.PullRequests[0]
	require.Equal(t, "release(development): v0.0.0-drc.0", *pr.Title)
	require.Equal(t, devReleaseBranch, *pr.Head.Ref)
	require.Equal(t, devBranch, *pr.Base.Ref)
	require.Equal(t, false, *pr.Draft)

	expected := []string{
		":robot: I have created a release candidate *beep* *boop*",
		"",
		"---",
		"",
		"## [v0.0.0-drc.0] Initial Version _2022-01-01 00:00 UTC_",
		"### Features",
		"",
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) init",
		"",
		"### Fixes",
		"",
		"- ([`hash3-h`](https://github.com/owner/name/commit/hash3-hash3)) another fix related to the feature",
		"- ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2)) fix related to the feature",
		"",
		"---",
		"",
		"This release was composed by [version_actions](https://github.com/jakbytes/version_actions)",
	}

	for i, line := range strings.Split(*pr.Body, "\n") {
		if strings.HasPrefix(line, "## [") {
			prefix := strings.Split(line, " (")[0]
			assert.True(t, strings.HasPrefix(expected[i], prefix))
		} else {
			assert.Equal(t, expected[i], line)
		}
	}
}

func testPrereleasePromotion(t *testing.T) {
	prs.PullRequests = []*github.PullRequest{} // pull request was merged into development
	repositories.Tags = []*github.RepositoryTag{
		{
			Name: github.String("v0.0.0-drc.0"),
		},
	} // no tags in the repository
	os.Args = []string{"program", "release", "token", repositoryOwner, repositoryName, devBranch, stagingBranch, stagingPrereleaseIdentifier, mainBranch, "none"}

	assert.NotPanics(t, func() {
		version.Execute()
	}, "promote generation should not panic")

	require.Equal(t, 1, len(prs.PullRequests))
	pr := prs.PullRequests[0]
	require.Equal(t, "release(staging): v0.0.0-src.0", *pr.Title)
	require.Equal(t, devPromoteBranch, *pr.Head.Ref)
	require.Equal(t, stagingBranch, *pr.Base.Ref)
	require.Equal(t, false, *pr.Draft)

	expected := []string{
		":robot: I have created a release candidate *beep* *boop*",
		"",
		"---",
		"",
		"## [v0.0.0-src.0] Initial Version _2022-01-01 00:00 UTC_",
		"### Features",
		"",
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) init",
		"",
		"### Fixes",
		"",
		"- ([`hash3-h`](https://github.com/owner/name/commit/hash3-hash3)) another fix related to the feature",
		"- ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2)) fix related to the feature",
		"",
		"---",
		"",
		"This release was composed by [version_actions](https://github.com/jakbytes/version_actions)",
	}

	for i, line := range strings.Split(*pr.Body, "\n") {
		if strings.HasPrefix(line, "## [") {
			prefix := strings.Split(line, " (")[0]
			assert.True(t, strings.HasPrefix(expected[i], prefix))
		} else {
			assert.Equal(t, expected[i], line)
		}
	}

	t.Run("Testing updating version release candidate", testUpdatePrereleasePromotion)

}

func testUpdatePrereleasePromotion(t *testing.T) {
	prs.PullRequests = []*github.PullRequest{} // pull request was merged into staging
	repositories.Tags = []*github.RepositoryTag{
		{
			Name: github.String("v0.0.0-src.0"),
			Commit: &github.Commit{
				SHA: github.String("hash3-hash3"),
			},
		},
	} // no tags in the repository
	repositories.Commits = []*github.RepositoryCommit{
		{
			SHA: github.String("hash4-hash4"),
			Commit: &github.Commit{
				Message: github.String("fix: yet another fix related to the feature"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
				},
				Tree: &github.Tree{
					SHA: github.String("tree1-tree1"),
				},
			},
		},
		{
			SHA: github.String("hash3-hash3"),
			Commit: &github.Commit{
				Message: github.String("fix: another fix related to the feature"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
				},
				Tree: &github.Tree{
					SHA: github.String("tree1-tree1"),
				},
			},
		},
		{
			SHA: github.String("hash2-hash2"),
			Commit: &github.Commit{
				Message: github.String("fix: fix related to the feature"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				},
				Tree: &github.Tree{
					SHA: github.String("tree1-tree1"),
				},
			},
		},
		{
			SHA: github.String("hash1-hash1"),
			Commit: &github.Commit{
				Message: github.String("feat: init"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
		},
	}

	repositories.Comparison.Commits = repositories.Commits

	os.Args = []string{"program", "release", "token", repositoryOwner, repositoryName, devBranch, stagingBranch, stagingPrereleaseIdentifier, mainBranch, "none"}

	assert.NotPanics(t, func() {
		version.Execute()
	}, "version generation should not panic")

	require.Equal(t, 1, len(prs.PullRequests))
	pr := prs.PullRequests[0]
	require.Equal(t, "release(staging): v0.0.0-src.1", *pr.Title)
	require.Equal(t, devPromoteBranch, *pr.Head.Ref)
	require.Equal(t, stagingBranch, *pr.Base.Ref)
	require.Equal(t, false, *pr.Draft)

	expected := []string{
		":robot: I have created a release candidate *beep* *boop*",
		"",
		"---",
		"",
		"## [v0.0.0-src.1] Initial Version _2022-01-01 00:00 UTC_",
		"### Features",
		"",
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) init",
		"",
		"### Fixes",
		"",
		"- ([`hash4-h`](https://github.com/owner/name/commit/hash4-hash4)) yet another fix related to the feature",
		"- ([`hash3-h`](https://github.com/owner/name/commit/hash3-hash3)) another fix related to the feature",
		"- ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2)) fix related to the feature",
		"",
		"---",
		"",
		"This release was composed by [version_actions](https://github.com/jakbytes/version_actions)",
	}

	for i, line := range strings.Split(*pr.Body, "\n") {
		if strings.HasPrefix(line, "## [") {
			prefix := strings.Split(line, " (")[0]
			assert.True(t, strings.HasPrefix(expected[i], prefix))
		} else {
			assert.Equal(t, expected[i], line)
		}
	}
}

func testPromotion(t *testing.T) {
	prs.PullRequests = []*github.PullRequest{} // pull request was merged into development
	repositories.Tags = []*github.RepositoryTag{
		{
			Name: github.String("v0.0.0-src.1"),
			Commit: &github.Commit{
				SHA: github.String("hash4-hash4"),
			},
		},
	} // no tags in the repository

	repositories.Commits = []*github.RepositoryCommit{
		{
			SHA: github.String("hash4-hash4"),
			Commit: &github.Commit{
				Message: github.String("fix: yet another fix related to the feature"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
				},
				Tree: &github.Tree{
					SHA: github.String("tree1-tree1"),
				},
			},
		},
		{
			SHA: github.String("hash3-hash3"),
			Commit: &github.Commit{
				Message: github.String("fix: another fix related to the feature"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
				},
				Tree: &github.Tree{
					SHA: github.String("tree1-tree1"),
				},
			},
		},
		{
			SHA: github.String("hash2-hash2"),
			Commit: &github.Commit{
				Message: github.String("fix: fix related to the feature"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				},
				Tree: &github.Tree{
					SHA: github.String("tree1-tree1"),
				},
			},
		},
		{
			SHA: github.String("hash1-hash1"),
			Commit: &github.Commit{
				Message: github.String("feat: init"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
		},
	}
	repositories.Comparison.Commits = repositories.Commits

	os.Args = []string{"program", "release", "token", repositoryOwner, repositoryName, stagingBranch, mainBranch, stagingPrereleaseIdentifier, mainBranch, "none"}

	assert.NotPanics(t, func() {
		version.Execute()
	}, "promote generation should not panic")

	require.Equal(t, 1, len(prs.PullRequests))
	pr := prs.PullRequests[0]
	require.Equal(t, "release(main): v0.0.0", *pr.Title)
	require.Equal(t, stagingPromoteBranch, *pr.Head.Ref)
	require.Equal(t, mainBranch, *pr.Base.Ref)
	require.Equal(t, false, *pr.Draft)

	expected := []string{
		":robot: I have created a release *beep* *boop*",
		"",
		"---",
		"",
		"## [v0.0.0] Initial Version _2022-01-01 00:00 UTC_",
		"### Features",
		"",
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) init",
		"",
		"### Fixes",
		"",
		"- ([`hash4-h`](https://github.com/owner/name/commit/hash4-hash4)) yet another fix related to the feature",
		"- ([`hash3-h`](https://github.com/owner/name/commit/hash3-hash3)) another fix related to the feature",
		"- ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2)) fix related to the feature",
		"",
		"---",
		"",
		"This release was composed by [version_actions](https://github.com/jakbytes/version_actions)",
	}

	for i, line := range strings.Split(*pr.Body, "\n") {
		if strings.HasPrefix(line, "## [") {
			prefix := strings.Split(line, " (")[0]
			assert.True(t, strings.HasPrefix(expected[i], prefix))
		} else {
			assert.Equal(t, expected[i], line)
		}
	}

	t.Run("Testing updating version promotion", testUpdatePromotion)

}

func testUpdatePromotion(t *testing.T) {
	repositories.Tags = []*github.RepositoryTag{
		{
			Name: github.String("v0.0.0-src.0"),
			Commit: &github.Commit{
				SHA: github.String("hash3-hash3"),
			},
		},
	} // no tags in the repository
	repositories.Commits = []*github.RepositoryCommit{
		{
			SHA: github.String("hash4-hash4"),
			Commit: &github.Commit{
				Message: github.String("fix: yet another fix related to the feature"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
				},
				Tree: &github.Tree{
					SHA: github.String("tree1-tree1"),
				},
			},
		},
		{
			SHA: github.String("hash3-hash3"),
			Commit: &github.Commit{
				Message: github.String("fix: another fix related to the feature"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
				},
				Tree: &github.Tree{
					SHA: github.String("tree1-tree1"),
				},
			},
		},
		{
			SHA: github.String("hash2-hash2"),
			Commit: &github.Commit{
				Message: github.String("fix: fix related to the feature"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				},
				Tree: &github.Tree{
					SHA: github.String("tree1-tree1"),
				},
			},
		},
		{
			SHA: github.String("hash1-hash1"),
			Commit: &github.Commit{
				Message: github.String("feat: init"),
				Author: &github.CommitAuthor{
					Login: github.String("login1"),
				},
				Committer: &github.CommitAuthor{
					Login: github.String("user1"),
					Date:  &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
		},
	}

	os.Args = []string{"program", "release", "token", repositoryOwner, repositoryName, stagingBranch, mainBranch, stagingPrereleaseIdentifier, mainBranch, "none"}

	assert.NotPanics(t, func() {
		version.Execute()
	}, "version generation should not panic")

	require.Equal(t, 1, len(prs.PullRequests))
	pr := prs.PullRequests[0]
	require.Equal(t, "release(main): v0.0.0", *pr.Title)
	require.Equal(t, stagingPromoteBranch, *pr.Head.Ref)
	require.Equal(t, mainBranch, *pr.Base.Ref)
	require.Equal(t, false, *pr.Draft)

	expected := []string{
		":robot: I have created a release *beep* *boop*",
		"",
		"---",
		"",
		"## [v0.0.0] Initial Version _2022-01-01 00:00 UTC_",
		"### Features",
		"",
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) init",
		"",
		"### Fixes",
		"",
		"- ([`hash4-h`](https://github.com/owner/name/commit/hash4-hash4)) yet another fix related to the feature",
		"- ([`hash3-h`](https://github.com/owner/name/commit/hash3-hash3)) another fix related to the feature",
		"- ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2)) fix related to the feature",
		"",
		"---",
		"",
		"This release was composed by [version_actions](https://github.com/jakbytes/version_actions)",
	}

	for i, line := range strings.Split(*pr.Body, "\n") {
		if strings.HasPrefix(line, "## [") {
			prefix := strings.Split(line, " (")[0]
			assert.True(t, strings.HasPrefix(expected[i], prefix))
		} else {
			assert.Equal(t, expected[i], line)
		}
	}
}
