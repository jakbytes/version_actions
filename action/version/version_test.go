package version

import (
	"context"
	"fmt"
	"github.com/jakbytes/version_actions/internal/mocks"
	"github.com/jakbytes/version_actions/tools/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSetup(t *testing.T) {
	os.Args = []string{"program", "version", "token", "owner", "name", "head", "base", "prereleaseIdentifier", "releaseBranch", "none"}
	client, args, err := setup()
	require.Nil(t, err)

	require.Equal(t, "token", args.Token)
	require.Equal(t, "owner", args.Owner)
	require.Equal(t, "name", args.Name)
	require.Equal(t, "head", args.Head)
	require.Equal(t, "base", args.Base)
	require.Equal(t, "prereleaseIdentifier", args.PrereleaseIdentifier)
	require.Equal(t, "releaseBranch", args.ReleaseBranch)
	require.NotNil(t, client)
}

func TestSetup_Panic(t *testing.T) {
	os.Args = []string{"program"}
	require.Panics(t, func() {
		_, _, _ = setup()
	})
}

func TestSetup_DefaultReleaseBranch(t *testing.T) {
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Repositories: &mocks.RepositoryService{},
		}
	}

	os.Args = []string{"program", "version", "token", "owner", "name", "head", "base", "prereleaseIdentifier", ".", "none"}
	client, args, err := setup()
	require.Nil(t, err)

	require.Equal(t, "token", args.Token)
	require.Equal(t, "owner", args.Owner)
	require.Equal(t, "name", args.Name)
	require.Equal(t, "head", args.Head)
	require.Equal(t, "base", args.Base)
	require.Equal(t, "prereleaseIdentifier", args.PrereleaseIdentifier)
	require.Equal(t, "main", args.ReleaseBranch)
	require.NotNil(t, client)
}

func TestSetup_DefaultReleaseBranch_Error(t *testing.T) {
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Repositories: &mocks.RepositoryService{
				GetError: assert.AnError,
			},
		}
	}

	os.Args = []string{"program", "version", "token", "owner", "name", "head", "base", "prereleaseIdentifier", ".", "none"}
	_, _, err := setup()
	require.NotNil(t, err)
	require.Equal(t, fmt.Errorf("failed to get default branch: %w", assert.AnError), err)
}

/*
func TestSetReleaseBranch_Create(t *testing.T) {
	count := 0
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base", "rc", "main"}
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Git: &mocks.GitService{},
			Repositories: &mocks.RepositoryService{
				GetBranchError: func(ctx context.Context, owner string, repo string, branch string, maxRedirects int) error {
					if branch == "release--branch--head" && count == 0 {
						count += 1
						return errors.New("404")
					}
					return nil
				},
			},
		}
	}
	client, args, err := setup()
	require.Nil(t, err)

	head, err := client.Repository().Branch(args.Head)
	require.Nil(t, err)

	branch, err := setReleaseBranch(args, client, head)
	require.Nil(t, err)

	require.Equal(t, "release--branch--head", branch.Name)
}

func TestSetReleaseBranch_Reset(t *testing.T) {
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base", "rc", "main"}
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Git:          &mocks.GitService{},
			Repositories: &mocks.RepositoryService{},
		}
	}
	client, args, err := setup()
	require.Nil(t, err)

	head, err := client.Repository().Branch(args.Head)
	require.Nil(t, err)

	branch, err := setReleaseBranch(args, client, head)
	require.Nil(t, err)

	require.Equal(t, "release--branch--head", branch.Name)
}

func TestSetReleaseBranch_Error(t *testing.T) {
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base", "rc", "main"}
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Git: &mocks.GitService{},
			Repositories: &mocks.RepositoryService{
				GetBranchError: func(ctx context.Context, owner string, repo string, branch string, maxRedirects int) error {
					if branch == "release--branch--head" {
						return assert.AnError
					}
					return nil
				},
			},
		}
	}
	client, args, err := setup()
	require.Nil(t, err)

	head, err := client.Repository().Branch(args.Head)
	require.Nil(t, err)

	_, err = setReleaseBranch(args, client, head)
	require.NotNil(t, err)
	require.Equal(t, assert.AnError, err)
}

func TestWriteChangelog(t *testing.T) {
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base", "rc", "head"}
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Repositories: &mocks.RepositoryService{},
		}
	}

	var exampleExisting = []string{
		"# Changelog",
		"",
		"---",
		"",
		"## [v1.1.0-beta.1]",
		"Feature 1",
		"",
		"---",
		"",
		"## [v1.0.0]",
		"Initial Version",
	}

	changelog.Path = "test_CHANGELOG.md"

	defer os.Remove(changelog.Path)

	err := changelog.WriteToFile(changelog.Path, exampleExisting)
	require.Nil(t, err)

	client, args, err := setup()
	require.Nil(t, err)

	head, err := client.Repository().Branch(args.Head)
	require.Nil(t, err)

	latest := &github.Version{
		Version: semver.MustParse("1.0.0"),
		RepositoryTag: &github.RepositoryTag{
			Commit: &github.Commit{
				SHA: github.String("hash1-hash1"),
			},
		},
	}

	nextVersion, cl, full, err := writeChangelog(args, client.Repository(), head, latest, nil)
	require.Nil(t, err)

	require.Equal(t, semver.MustParse("1.1.0").String(), nextVersion.String())

	assert.Equal(t, 9, len(cl))

	changelogLines := cl
	changelogLines = append([]string{"# Changelog", "", "---", ""}, changelogLines...)
	changelogLines = append(changelogLines, "---", "", "## [v1.0.0]", "Initial Version")

	expected := []string{
		"# Changelog",
		"",
		"---",
		"",
		"## [v1.1.0](https://github.com/owner/name/compare/v1.0.0...v1.1.0) _2024-01-31 18:13 UTC_",
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
		"## [v1.0.0]",
		"Initial Version",
	}

	i := 0
	err = utility.Open(changelog.Path, func(file *os.File) error {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "## [") {
				prefix := strings.Split(line, " _")[0]
				assert.True(t, strings.HasPrefix(changelogLines[i], prefix))
				assert.True(t, strings.HasPrefix(expected[i], prefix))
			} else {
				assert.Equal(t, changelogLines[i], line)
				assert.Equal(t, expected[i], line)
			}
			assert.Equal(t, full[i], line)
			i += 1
		}
		return nil
	})
	require.Nil(t, err)
}

func TestWriteChangelog_NoIncrement(t *testing.T) {
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "head", "rc", "head"}
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Repositories: &mocks.RepositoryService{
				Tags: []*github.RepositoryTag{
					{
						Name: github.String("v1.1.0"),
						Commit: &github.Commit{
							SHA: github.String("hash1-hash1"),
						},
					},
					{
						Name: github.String("v1.0.0"),
						Commit: &github.Commit{
							SHA: github.String("hash0-hash0"),
						},
					},
				},
				Commits: []*github.RepositoryCommit{
					{
						SHA: github.String("hash3-hash3"),
						Commit: &github.Commit{
							Message: github.String("docs: message1"),
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
						SHA: github.String("hash1-hash1"),
						Commit: &github.Commit{
							Message: github.String("feat: message1"),
							Author: &github.CommitAuthor{
								Login: github.String("login1"),
							},
							Committer: &github.CommitAuthor{
								Login: github.String("user1"),
								Date:  &github.Timestamp{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
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
								Date:  &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
							},
						},
					},
					{
						SHA: github.String("hash0-hash0"),
						Commit: &github.Commit{
							Message: github.String("Initial"),
							Author: &github.CommitAuthor{
								Login: github.String("login2"),
							},
							Committer: &github.CommitAuthor{
								Login: github.String("user2"),
								Date:  &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
							},
						},
					},
				},
			},
		}
	}

	var exampleExisting = []string{
		"# Changelog",
		"",
		"---",
		"",
		"## [v1.1.0](https://github.com/owner/name/compare/v1.0.0...v1.1.0) _2024-01-31 18:13 UTC_",
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
		"## [v1.0.0]",
		"Initial Version",
	}

	changelog.Path = "test_CHANGELOG.md"
	defer os.Remove(changelog.Path)

	err := changelog.WriteToFile(changelog.Path, exampleExisting)
	require.Nil(t, err)

	client, args, err := setup()
	require.Nil(t, err)

	head, err := client.Repository().Branch(args.Head)
	require.Nil(t, err)

	latest := &github.Version{
		Version: semver.MustParse("1.1.0"),
		RepositoryTag: &github.RepositoryTag{
			Commit: &github.Commit{
				SHA: github.String("hash1-hash1"),
			},
		},
	}

	nextVersion, _, _, err := writeChangelog(args, client.Repository(), head, latest, nil)
	require.Nil(t, err)

	require.Equal(t, semver.MustParse("1.1.0").String(), nextVersion.String())

	expected := []string{
		"# Changelog",
		"",
		"---",
		"",
		"## [v1.1.0](https://github.com/owner/name/compare/v1.0.0...v1.1.0) _2024-01-31 18:13 UTC_",
		"### Features",
		"",
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) message1",
		"",
		"### Fixes",
		"",
		"- ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2)) message2",
		"",
		"### Documentation",
		"",
		"- ([`hash3-h`](https://github.com/owner/name/commit/hash3-hash3)) message1",
		"",
		"---",
		"",
		"## [v1.0.0]",
		"Initial Version",
	}

	i := 0
	err = utility.Open(changelog.Path, func(file *os.File) error {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "## [") {
				prefix := strings.Split(line, " _")[0]
				assert.True(t, strings.HasPrefix(expected[i], prefix))
			} else {
				assert.Equal(t, expected[i], line)
			}
			i += 1
		}
		return nil
	})
	require.Nil(t, err)
}

func TestWriteChangelog_NoIncrement_RC(t *testing.T) {
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base", "rc", "base"}
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Repositories: &mocks.RepositoryService{
				Tags: []*github.RepositoryTag{
					{
						Name: github.String("v1.1.0-rc.0"),
						Commit: &github.Commit{
							SHA: github.String("hash1-hash1"),
						},
					},
					{
						Name: github.String("v1.0.0"),
						Commit: &github.Commit{
							SHA: github.String("hash0-hash0"),
						},
					},
				},
				Comparison: &github.CommitsComparison{
					Commits: []*github.RepositoryCommit{
						{
							SHA: github.String("hash3-hash3"),
							Commit: &github.Commit{
								Message: github.String("docs: message1"),
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
							SHA: github.String("hash1-hash1"),
							Commit: &github.Commit{
								Message: github.String("feat: message1"),
								Author: &github.CommitAuthor{
									Login: github.String("login1"),
								},
								Committer: &github.CommitAuthor{
									Login: github.String("user1"),
									Date:  &github.Timestamp{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
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
									Date:  &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
								},
							},
						},
						{
							SHA: github.String("hash0-hash0"),
							Commit: &github.Commit{
								Message: github.String("Initial"),
								Author: &github.CommitAuthor{
									Login: github.String("login2"),
								},
								Committer: &github.CommitAuthor{
									Login: github.String("user2"),
									Date:  &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
								},
							},
						},
					},
				},
				Commits: []*github.RepositoryCommit{
					{
						SHA: github.String("hash3-hash3"),
						Commit: &github.Commit{
							Message: github.String("docs: message1"),
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
						SHA: github.String("hash1-hash1"),
						Commit: &github.Commit{
							Message: github.String("feat: message1"),
							Author: &github.CommitAuthor{
								Login: github.String("login1"),
							},
							Committer: &github.CommitAuthor{
								Login: github.String("user1"),
								Date:  &github.Timestamp{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
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
								Date:  &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
							},
						},
					},
					{
						SHA: github.String("hash0-hash0"),
						Commit: &github.Commit{
							Message: github.String("Initial"),
							Author: &github.CommitAuthor{
								Login: github.String("login2"),
							},
							Committer: &github.CommitAuthor{
								Login: github.String("user2"),
								Date:  &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
							},
						},
					},
				},
			},
		}
	}

	var exampleExisting = []string{
		"# Changelog",
		"",
		"---",
		"",
		"## [v1.1.0-rc.0](https://github.com/owner/name/compare/v1.0.0...v1.1.0) (2024-01-31 18:13 UTC)",
		"### Features",
		"",
		"- message1",
		"  > _Contributed by [](https://github.com/) on 2024-01-31 13:13 UTC_ ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1))",
		"",
		"### Fixes",
		"",
		"- message2",
		"  > _Contributed by [](https://github.com/) on 2024-01-31 13:13 UTC_ ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2))",
		"",
		"---",
		"",
		"## [v1.0.0]",
		"Initial Version",
	}

	changelog.Path = "test_CHANGELOG.md"

	defer os.Remove(changelog.Path)

	err := changelog.WriteToFile(changelog.Path, exampleExisting)
	require.Nil(t, err)

	client, args, err := setup()
	require.Nil(t, err)

	head, err := client.Repository().Branch(args.Head)
	require.Nil(t, err)

	latest := &github.Version{
		Version: semver.MustParse("1.0.0"),
		RepositoryTag: &github.RepositoryTag{
			Commit: &github.Commit{
				SHA: github.String("hash0-hash0"),
			},
		},
	}

	prerelease := &github.Version{
		Version: semver.MustParse("1.1.0-rc.0"),
		RepositoryTag: &github.RepositoryTag{
			Commit: &github.Commit{
				SHA: github.String("hash1-hash1"),
			},
		},
	}

	nextVersion, _, _, err := writeChangelog(args, client.Repository(), head, latest, prerelease)
	require.Nil(t, err)

	require.Equal(t, semver.MustParse("1.1.0-rc.1").String(), nextVersion.String())

	expected := []string{
		"# Changelog",
		"",
		"---",
		"",
		"## [v1.1.0-rc.1](https://github.com/owner/name/compare/v1.0.0...v1.1.0-rc.1) _2024-01-31 18:13 UTC_",
		"### Features",
		"",
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) message1",
		"",
		"### Fixes",
		"",
		"- ([`hash2-h`](https://github.com/owner/name/commit/hash2-hash2)) message2",
		"",
		"### Documentation",
		"",
		"- ([`hash3-h`](https://github.com/owner/name/commit/hash3-hash3)) message1",
		"",
		"---",
		"",
		"## [v1.0.0]",
		"Initial Version",
	}

	i := 0
	err = utility.Open(changelog.Path, func(file *os.File) error {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "## [") {
				prefix := strings.Split(line, " _")[0]
				assert.True(t, strings.HasPrefix(expected[i], prefix))
			} else {
				assert.Equal(t, expected[i], line)
			}
			i += 1
		}
		return nil
	})
	require.Nil(t, err)
}

func TestComposePullRequest(t *testing.T) {
	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base", "rc", "head"}
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		return &github.Client{
			Repositories: &mocks.RepositoryService{},
		}
	}

	var exampleExisting = []string{
		"# Changelog",
		"",
		"---",
		"",
		"## [v1.1.0-beta.1]",
		"Feature 1",
		"",
		"---",
		"",
		"## [v1.0.0]",
		"Initial Version",
	}

	changelog.Path = "test_CHANGELOG.md"

	defer os.Remove(changelog.Path)

	err := changelog.WriteToFile(changelog.Path, exampleExisting)
	require.Nil(t, err)

	client, args, err := setup()
	require.Nil(t, err)

	head, err := client.Repository().Branch(args.Head)
	require.Nil(t, err)

	latest := &github.Version{
		Version: semver.MustParse("1.0.0"),
		RepositoryTag: &github.RepositoryTag{
			Commit: &github.Commit{
				SHA: github.String("hash1-hash1"),
			},
		},
	}

	nextVersion, cl, _, err := writeChangelog(args, client.Repository(), head, latest, nil)
	require.Nil(t, err)

	require.Equal(t, semver.MustParse("1.1.0").String(), nextVersion.String())

	assert.Equal(t, 9, len(cl))

	changelogLines := cl
	changelogLines = append([]string{"# Changelog", "", "---", ""}, changelogLines...)
	changelogLines = append(changelogLines, "---", "", "## [v1.0.0]", "Initial Version")

	expected := []string{
		"# Changelog",
		"",
		"---",
		"",
		"## [v1.1.0](https://github.com/owner/name/compare/v1.0.0...v1.1.0) _2024-01-31 18:13 UTC_",
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
		"## [v1.0.0]",
		"Initial Version",
	}

	i := 0
	err = utility.Open(changelog.Path, func(file *os.File) error {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "## [") {
				prefix := strings.Split(line, " _")[0]
				assert.True(t, strings.HasPrefix(changelogLines[i], prefix))
				assert.True(t, strings.HasPrefix(expected[i], prefix))
			} else {
				assert.Equal(t, changelogLines[i], line)
				assert.Equal(t, expected[i], line)
			}
			i += 1
		}
		return nil
	})
	require.Nil(t, err)

	title, body := composePullRequest(args, nextVersion, cl, true)
	require.Equal(t, "release(head): v1.1.0", title)

	expected = []string{
		":robot: I have created a release *beep* *boop*",
		"",
		"---",
		"",
		"## [v1.1.0](https://github.com/owner/name/compare/v1.0.0...v1.1.0) _2024-01-31 18:13 UTC_",
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
		"This PR was generated by [semver-action.](https://github.com/jacobkring/semver-action)",
	}

	for i, line := range body {
		if strings.HasPrefix(line, "## [") {
			prefix := strings.Split(line, " _")[0]
			assert.True(t, strings.HasPrefix(expected[i], prefix))
		} else {
			assert.Equal(t, expected[i], line)
		}
	}
}

func TestVersion(t *testing.T) {
	changelog.Path = "test_CHANGELOG.md"

	defer os.Remove(changelog.Path)

	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base", "rc", "head"}
	var client *github.Client

	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		client = &github.Client{
			Repositories: &mocks.RepositoryService{},
			Git:          &mocks.GitService{},
			PullRequests: &mocks.PullRequestsService{},
			RepositoryMetadata: github.RepositoryMetadata{
				Owner: owner,
				Name:  name,
			},
		}
		return client
	}

	err := version()
	require.Nil(t, err)

	assert.Equal(t, 1, len(client.PullRequests.(*mocks.PullRequestsService).PullRequests))
}

func TestVersion_NoTags(t *testing.T) {
	changelog.Path = "test_CHANGELOG.md"

	defer os.Remove(changelog.Path)

	os.Args = []string{"program", "action", "token", "owner", "name", "head", "base", "rc", "base"}
	var client *github.Client
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		client = &github.Client{
			Repositories: &mocks.RepositoryService{
				Tags: []*github.RepositoryTag{},
			},
			Git:          &mocks.GitService{},
			PullRequests: &mocks.PullRequestsService{},
			RepositoryMetadata: github.RepositoryMetadata{
				Owner: owner,
				Name:  name,
			},
		}
		return client
	}

	err := version()
	require.Nil(t, err)

	assert.Equal(t, 1, len(client.PullRequests.(*mocks.PullRequestsService).PullRequests))
	pr := client.PullRequests.(*mocks.PullRequestsService).PullRequests[0]
	pr.Base.Ref = github.String("base")
	pr.Head.Ref = github.String("head")
}

func TestVersion_ReleaseBranch(t *testing.T) {
	changelog.Path = "test_CHANGELOG.md"

	defer os.Remove(changelog.Path)

	os.Args = []string{"program", "action", "token", "owner", "name", "head", "head", "rc", "head"}
	var client *github.Client
	NewClient = func(ctx context.Context, token string, owner string, name string) *github.Client {
		client = &github.Client{
			Repositories: &mocks.RepositoryService{},
			Git:          &mocks.GitService{},
			PullRequests: &mocks.PullRequestsService{},
			RepositoryMetadata: github.RepositoryMetadata{
				Owner: owner,
				Name:  name,
			},
		}
		return client
	}

	err := version()
	require.Nil(t, err)

	assert.Equal(t, 1, len(client.PullRequests.(*mocks.PullRequestsService).PullRequests))

	pr := client.PullRequests.(*mocks.PullRequestsService).PullRequests[0]
	pr.Base.Ref = github.String("head")
	pr.Head.Ref = github.String("release--branch--head")

	body := strings.Split(*pr.Body, "\n")
	expected := []string{
		":robot: I have created a release *beep* *boop*",
		"",
		"---",
		"",
		"## [v1.1.0](https://github.com/owner/name/compare/v1.0.1...v1.1.0) _2024-01-31 18:13 UTC_",
		"### Features",
		"",
		"- ([`hash1-h`](https://github.com/owner/name/commit/hash1-hash1)) message1",
		"",
		"---",
		"",
		"This PR was generated by [semver-action.](https://github.com/jacobkring/semver-action)",
	}

	for i, line := range body {
		if strings.HasPrefix(line, "## [") {
			prefix := strings.Split(line, " _")[0]
			assert.True(t, strings.HasPrefix(expected[i], prefix))
		} else {
			assert.Equal(t, expected[i], line)
		}
	}

}

*/
