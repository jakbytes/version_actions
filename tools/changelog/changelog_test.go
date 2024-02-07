package changelog

import (
	"bufio"
	"os"
	"strings"
	"testing"
	"time"
	"version_actions/internal/utility"
	"version_actions/tools/conventional"
	"version_actions/tools/semver"

	"github.com/google/go-github/v58/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateVersionHeader(t *testing.T) {
	org := "exampleOrg"
	repo := "exampleRepo"
	version, _ := semver.NewVersion("1.0.0")
	prevVersion, _ := semver.NewVersion("0.9.0")

	// Test with disableVersionHeader = true
	result := generateVersionHeader(org, repo, prevVersion, version, true)
	assert.Equal(t, "## Changelog", result)

	// Test with previousVersion = nil
	result = generateVersionHeader(org, repo, nil, version, false)
	assert.Contains(t, result, "## [v1.0.0] Initial Version", "Header should contain initial version info")

	// Test with previousVersion != nil
	result = generateVersionHeader(org, repo, prevVersion, version, false)
	assert.Contains(t, result, "https://github.com/exampleOrg/exampleRepo/compare/v0.9.0...v1.0.0", "Header should contain version comparison link")
}

func TestFormatCommit(t *testing.T) {
	// Mock commit data
	date := time.Date(2021, 10, 1, 12, 0, 0, 0, time.UTC)
	commit := &github.RepositoryCommit{
		Commit: &github.Commit{
			Message: github.String("Test commit: This is a test"),
			Committer: &github.CommitAuthor{
				Date: &github.Timestamp{Time: date},
				Name: github.String("John Doe"),
			},
		},
		SHA:    github.String("1234567890abcdef"),
		Author: &github.User{Login: github.String("johndoe")},
	}

	// Expected format
	expected := "- ([`1234567`](https://github.com/org/repo/commit/1234567890abcdef)) This is a test"

	// Running the test with assert
	result := formatCommit("org", "repo", commit)
	assert.Equal(t, Markdown(strings.Split(expected, "\n")), result, "formatCommit should format the commit correctly")
}

// Sample mock commit function
func mockCommit(message, authorName, authorLogin, sha string) *github.RepositoryCommit {
	date := time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)
	return &github.RepositoryCommit{
		Commit: &github.Commit{
			Message: github.String(message),
			Committer: &github.CommitAuthor{
				Name: github.String(authorName),
				Date: &github.Timestamp{Time: date},
			},
		},
		SHA: github.String(sha),
		Author: &github.User{
			Login: github.String(authorLogin),
		},
	}
}

func TestGenerateNewChangelog(t *testing.T) {
	org := "exampleOrg"
	repo := "exampleRepo"
	version, _ := semver.NewVersion("1.0.0")

	// Mock commits for each category
	breaking := []*github.RepositoryCommit{mockCommit("feat!: breaking change", "Alice", "alice", "abc1234")}
	feat := []*github.RepositoryCommit{mockCommit("feat: new feature", "Bob", "bob", "def5678")}
	fix := []*github.RepositoryCommit{mockCommit("fix: bug fix", "Charlie", "charlie", "ghi9012")}

	// Test with non-empty commit lists
	changelog := GenerateNewChangelog(org, repo, nil, version, conventional.Commits{Breaking: breaking, Feat: feat, Fix: fix}, false)

	require.Equal(t, 13, len(changelog))
	require.True(t, strings.HasPrefix(changelog[0], "## [v1.0.0]"), "Changelog should contain version header")

	// Test with empty commit lists and disableVersionHeader = true
	changelog = GenerateNewChangelog(org, repo, nil, version, conventional.Commits{}, true)
	assert.NotContains(t, changelog, "## v1.0.0", "Changelog should not contain version header when disabled")
	assert.NotContains(t, changelog, "Breaking Changes", "Changelog should not contain Breaking Changes section for empty list")
}

func TestShouldSkipLine_Header(t *testing.T) {
	line := "# Changelog"
	currentVersion := false
	skipNextBreak := false
	skipNextSpace := false
	versionHeading := "## [v1.0.0"
	assert.True(t, shouldSkipLine(line, &currentVersion, &skipNextBreak, &skipNextSpace, versionHeading), "shouldSkipLine should return true for changelog header")
	assert.True(t, skipNextBreak, "shouldSkipLine should set skipNextBreak to true for changelog header")
}

func TestShouldSkipLine_CurrentVersion(t *testing.T) {
	line := "## [v1.0.0]("
	currentVersion := false
	skipNextBreak := false
	skipNextSpace := false
	versionHeading := "## [v1.0.0"
	assert.True(t, shouldSkipLine(line, &currentVersion, &skipNextBreak, &skipNextSpace, versionHeading), "shouldSkipLine should return true for current version header")
	assert.True(t, currentVersion, "shouldSkipLine should set currentVersion to true for current version header")

	line = "---"
	currentVersion = true
	skipNextBreak = false
	versionHeading = "## [v1.0.0"
	// should not skip line on different version
	assert.False(t, shouldSkipLine(line, &currentVersion, &skipNextBreak, &skipNextSpace, versionHeading), "shouldSkipLine should return false for different version header")
	assert.False(t, currentVersion, "shouldSkipLine should not change currentVersion for different version header")
}

func TestShouldSkipLine_SkipNextBreak(t *testing.T) {
	line := "---"
	currentVersion := false
	skipNextBreak := true
	skipNextSpace := false
	versionHeading := "## [v1.0.0"
	assert.True(t, shouldSkipLine(line, &currentVersion, &skipNextBreak, &skipNextSpace, versionHeading), "shouldSkipLine should return true for skipped break")
	assert.False(t, skipNextBreak, "shouldSkipLine should set skipNextBreak to false for skipped break")
}

func TestUpdateChangelog(t *testing.T) {
	var testInput = []string{
		"# Changelog",
		"---",
		"## [v1.1.0-beta.2]",
		"Feature 2",
	}

	var exampleExisting = []string{
		"# Changelog",
		"---",
		"## [v1.1.0-beta.1]",
		"Feature 1",
		"---",
		"## [v1.0.0]",
		"Initial Version",
	}

	Path = "test_CHANGELOG.md"

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(err)
		}
	}(Path)

	err := WriteToFile(Path, exampleExisting)
	require.Nil(t, err)

	version, err := semver.NewVersion("1.1.0-beta.2")
	require.Nil(t, err)

	lines, err := updateChangelog(version, testInput)
	require.Nil(t, err)

	assert.Equal(t, 7, len(lines), "updateChangelog should return 5 lines")
	assert.Equal(t, "# Changelog", lines[0], "updateChangelog should retain header")
	assert.Equal(t, "---", lines[1], "updateChangelog should retain break")
	assert.Equal(t, "## [v1.1.0-beta.2]", lines[2], "updateChangelog should update version header")
	assert.Equal(t, "Feature 2", lines[3], "updateChangelog should retain feature")
	assert.Equal(t, "---", lines[4], "updateChangelog should retain break")
	assert.Equal(t, "## [v1.0.0]", lines[5], "updateChangelog should retain previous version header")
	assert.Equal(t, "Initial Version", lines[6], "updateChangelog should retain previous version feature")
}

func TestWriteChangelog(t *testing.T) {
	var exampleExisting = []string{
		"# Changelog",
		"---",
		"## [v1.1.0-beta.1]",
		"Feature 1",
		"---",
		"## [v1.0.0]",
		"Initial Version",
	}

	Path = "test_CHANGELOG.md"

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(err)
		}
	}(Path)

	err := WriteToFile(Path, exampleExisting)
	require.Nil(t, err)

	org := "exampleOrg"
	repo := "exampleRepo"
	version, _ := semver.NewVersion("1.1.0")
	prevVersion, _ := semver.NewVersion("1.0.0")

	// Mock commits for each category
	breaking := []*github.RepositoryCommit{mockCommit("feat!: breaking change", "Alice", "alice", "abc1234")}
	feat := []*github.RepositoryCommit{mockCommit("feat: new feature", "Bob", "bob", "def5678")}
	fix := []*github.RepositoryCommit{mockCommit("fix: bug fix", "Charlie", "charlie", "ghi9012")}

	changelog, _, err := WriteChangelog(org, repo, prevVersion, version, conventional.Commits{Breaking: breaking, Feat: feat, Fix: fix}, false)
	require.Nil(t, err)

	assert.Equal(t, 13, len(changelog), "WriteChangelog should return 16 lines")

	changelogLines := append([]string{"# Changelog", "", "---", ""}, changelog...)
	changelogLines = append(changelogLines, "---", "## [v1.0.0]", "Initial Version")

	i := 0
	err = utility.Open(Path, func(file *os.File) error {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			assert.Equal(t, changelogLines[i], line)
			i += 1
		}
		return nil
	})
	require.Nil(t, err)
}

func TestWriteChangelog_UpdateChangelogError(t *testing.T) {
	_, _ = os.Create("test_CHANGELOG.md")
	Path = "test_CHANGELOG.md"
	UpdateChangelog = func(version *semver.Version, lines Markdown) (Markdown, error) {
		return nil, assert.AnError
	}
	defer func() {
		UpdateChangelog = updateChangelog
	}()

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(err)
		}
	}(Path)

	_, err := UpdateChangelog(nil, nil)
	require.Equal(t, assert.AnError, err)

	org := "exampleOrg"
	repo := "exampleRepo"
	version, _ := semver.NewVersion("1.1.0")
	prevVersion, _ := semver.NewVersion("1.0.0")

	// Mock commits for each category
	breaking := []*github.RepositoryCommit{mockCommit("feat!: breaking change", "Alice", "alice", "abc1234")}
	feat := []*github.RepositoryCommit{mockCommit("feat: new feature", "Bob", "bob", "def5678")}
	fix := []*github.RepositoryCommit{mockCommit("fix: bug fix", "Charlie", "charlie", "ghi9012")}

	_, _, err = WriteChangelog(org, repo, prevVersion, version, conventional.Commits{Breaking: breaking, Feat: feat, Fix: fix}, false)
	require.NotNil(t, err)
	require.Equal(t, assert.AnError, err)
}

func TestWriteChangelog_WriteToFileError(t *testing.T) {
	_, _ = os.Create("test_CHANGELOG.md")
	Path = "test_CHANGELOG.md"
	WriteString = func(file *os.File, line string) error {
		return assert.AnError
	}
	defer func() {
		WriteString = writeString
	}()

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(err)
		}
	}(Path)

	org := "exampleOrg"
	repo := "exampleRepo"
	version, _ := semver.NewVersion("1.1.0")
	prevVersion, _ := semver.NewVersion("1.0.0")

	// Mock commits for each category
	breaking := []*github.RepositoryCommit{mockCommit("feat!: breaking change", "Alice", "alice", "abc1234")}
	feat := []*github.RepositoryCommit{mockCommit("feat: new feature", "Bob", "bob", "def5678")}
	fix := []*github.RepositoryCommit{mockCommit("fix: bug fix", "Charlie", "charlie", "ghi9012")}

	_, _, err := WriteChangelog(org, repo, prevVersion, version, conventional.Commits{Breaking: breaking, Feat: feat, Fix: fix}, false)
	require.NotNil(t, err)
	require.Equal(t, assert.AnError.Error(), err.Error())
}

func TestShouldSkipLine_SkipNextSpace(t *testing.T) {
	line := ""
	currentVersion := false
	skipNextBreak := false
	skipNextSpace := true
	versionHeading := "## [v1.0.0"
	assert.True(t, shouldSkipLine(line, &currentVersion, &skipNextBreak, &skipNextSpace, versionHeading), "shouldSkipLine should return true for skipped space")
	assert.False(t, skipNextSpace, "shouldSkipLine should set skipNextSpace to false for skipped space")
}
