package changelog

import (
	"bufio"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"testing"
	"time"
	"version_actions/internal/logger"
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
	log.Logger = logger.Base()
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

	commit = &github.RepositoryCommit{
		Commit: &github.Commit{
			Message: github.String(strings.Join([]string{"feat(development): a solid description",
				"",
				"a ton more detail to really understand what's going on",
				"",
				"Introduce a request id and a reference to latest request. Dismiss",
				"incoming responses other than from latest request.",
				"",
				"Remove timeouts which were used to mitigate the racing issue but are",
				"obsolete now.",
				"",
				"BREAKING CHANGE: use JavaScript features not available in Node 6."}, "\n")),
			Committer: &github.CommitAuthor{
				Date: &github.Timestamp{Time: date},
				Name: github.String("John Doe"),
			},
		},
		SHA:    github.String("1234567890abcdef"),
		Author: &github.User{Login: github.String("johndoe")},
	}

	// Expected format
	e := Markdown{
		"- ([`1234567`](https://github.com/org/repo/commit/1234567890abcdef)) a solid description",
		"  > ",
		"  > a ton more detail to really understand what's going on",
		"  > ",
		"  > Introduce a request id and a reference to latest request. Dismiss",
		"  > incoming responses other than from latest request.",
		"  > ",
		"  > Remove timeouts which were used to mitigate the racing issue but are",
		"  > obsolete now.",
		"  > ",
		"  > BREAKING CHANGE: use JavaScript features not available in Node 6.",
	}

	// Running the test with assert
	result = formatCommit("org", "repo", commit)
	for i, line := range e {
		assert.Equal(t, line, result[i])
	}
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
	assert.True(t, skipLine(line, &currentVersion, &skipNextBreak, &skipNextSpace, versionHeading), "skipLine should return true for changelog header")
	assert.True(t, skipNextBreak, "skipLine should set skipNextBreak to true for changelog header")
}

func TestShouldSkipLine_CurrentVersion(t *testing.T) {
	line := "## [v1.0.0]("
	currentVersion := false
	skipNextBreak := false
	skipNextSpace := false
	versionHeading := "## [v1.0.0"
	assert.True(t, skipLine(line, &currentVersion, &skipNextBreak, &skipNextSpace, versionHeading), "skipLine should return true for current version header")
	assert.True(t, currentVersion, "skipLine should set currentVersion to true for current version header")

	line = "## [v0.9.0]("
	currentVersion = true
	skipNextBreak = false
	versionHeading = "## [v1.0.0"
	// should not skip line on different version
	assert.False(t, skipLine(line, &currentVersion, &skipNextBreak, &skipNextSpace, versionHeading), "skipLine should return false for different version header")
	assert.False(t, currentVersion, "skipLine should not change currentVersion for different version header")
}

func TestShouldSkipLine_SkipNextBreak(t *testing.T) {
	line := "---"
	currentVersion := false
	skipNextBreak := true
	skipNextSpace := false
	versionHeading := "## [v1.0.0"
	assert.True(t, skipLine(line, &currentVersion, &skipNextBreak, &skipNextSpace, versionHeading), "skipLine should return true for skipped break")
	assert.False(t, skipNextBreak, "skipLine should set skipNextBreak to false for skipped break")
}

func TestUpdateChangelog(t *testing.T) {
	var testInput = []string{
		"# Changelog",
		"",
		"## [v1.1.0-beta.2]",
		"Feature 2",
		"",
	}

	var exampleExisting = []string{
		"# Changelog",
		"",
		"## [v1.1.0-beta.1]",
		"Feature 1",
		"",
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

	assert.Equal(t, 7, len(lines))
	assert.Equal(t, "# Changelog", lines[0])
	assert.Equal(t, "", lines[1])
	assert.Equal(t, "## [v1.1.0-beta.2]", lines[2])
	assert.Equal(t, "Feature 2", lines[3])
	assert.Equal(t, "", lines[4])
	assert.Equal(t, "## [v1.0.0]", lines[5])
	assert.Equal(t, "Initial Version", lines[6])
}

func TestWriteChangelog(t *testing.T) {
	var exampleExisting = []string{
		"# Changelog",
		"",
		"## [v1.1.0-beta.1]",
		"Feature 1",
		"",
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

	changelogLines := append([]string{"# Changelog", ""}, changelog...)
	changelogLines = append(changelogLines, "## [v1.0.0]", "Initial Version")

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
	assert.True(t, skipLine(line, &currentVersion, &skipNextBreak, &skipNextSpace, versionHeading), "skipLine should return true for skipped space")
	assert.False(t, skipNextSpace, "skipLine should set skipNextSpace to false for skipped space")
}

func TestUpdateChangelog_Scenario1(t *testing.T) {
	UpdateChangelog = updateChangelog
	log.Logger = logger.Base()
	var testInput = []string{
		"# Changelog",
		"",
		"## [v0.1.0-drc.0](https://github.com/jakbytes/version_actions/compare/v0.0.0...v0.1.0-drc.0) (2024-02-07)",
		"### Features",
		"",
		"- ([`38f1bd1`](https://github.com/jakbytes/version_actions/commit/38f1bd1091e162416bbcc653da5865b8f70e2c49)) breaking changes text capitalized to call it out strongly",
		"- ([`7237226`](https://github.com/jakbytes/version_actions/commit/72372265d197605918b127c92eb75375c3715382)) date on version is simplified",
		"- ([`0ba489f`](https://github.com/jakbytes/version_actions/commit/0ba489f5f33d221061c149fed64166c26c6322ae)) extract prerelease identifier action",
		"",
		"## [v0.0.0] Initial Version (2024-02-07)",
		"### Features",
		"",
		"- ([`c0d1dcd`](https://github.com/jakbytes/version_actions/commit/c0d1dcd0e3483390d8d7405569bcf3eadcce5710)) initial supported actions, version, sync, pull_request, extract_commit, download_release_asset",
		"",
		"### Fixes",
		"",
		"- ([`58bf05c`](https://github.com/jakbytes/version_actions/commit/58bf05caf571984ec6b2233ddb6f18a109a624ba)) type value needs to be output for further activity",
		"- ([`e1729a9`](https://github.com/jakbytes/version_actions/commit/e1729a947a61a321155939e72779334c88033b47)) action trigger should be set properly",
		"- ([`ba3d06f`](https://github.com/jakbytes/version_actions/commit/ba3d06fc58c65dc4fae5dd39c0d539207d906118)) hanging % needed to be removed from version action",
		"- ([`19bfb4d`](https://github.com/jakbytes/version_actions/commit/19bfb4db2aa5af63bead5067d2d3582e6b67fba2)) don't use best effort",
		"- ([`1a481d7`](https://github.com/jakbytes/version_actions/commit/1a481d72d0715ae6d7d88a9b434502513529c18c)) should be using v4 actions checkout",
		"- ([`1487ff3`](https://github.com/jakbytes/version_actions/commit/1487ff34f740541c9cb5aa3345aa14e6d1d93abc)) commits should be freeform to allow release and others",
		"- ([`68906c8`](https://github.com/jakbytes/version_actions/commit/68906c816d30d62c6f67c4a35b5e6003ccd74fbf)) download_release_asset shouldnt have quotes around the chmod val, version should not modify yml",
		"- ([`42328c0`](https://github.com/jakbytes/version_actions/commit/42328c0dc7d95b59e58c1373f678834420f8c329)) actions should reference version_action, not action",
		"- ([`8d24825`](https://github.com/jakbytes/version_actions/commit/8d24825ef39953f45c2fae275b420777c635ba5c)) a few more references to the old path were not adjusted",
		"- ([`db31802`](https://github.com/jakbytes/version_actions/commit/db31802dc409e7306ca2a4b17a8a1ba3e8332c05)) use the download_release_asset in pull_request, rename action.go to version_action.go",
		"",
		"### CI/CD",
		"",
		"- ([`a48f0ae`](https://github.com/jakbytes/version_actions/commit/a48f0aeac3a5c4ce3bed5af4e055bff7174bd99f)) fix reference to type",
		"- ([`3b55e7f`](https://github.com/jakbytes/version_actions/commit/3b55e7fbce860c789836006c2c1e93ab3a1554ce)) actions need to reference the correct path",
		"- ([`ed5f7a3`](https://github.com/jakbytes/version_actions/commit/ed5f7a398dd060d3a9769c344206c2b86dad2959)) remove debugging action",
	}

	Path = "test_CHANGELOG.md"

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(err)
		}
	}(Path)

	err := WriteToFile(Path, testInput)
	require.Nil(t, err)

	feat := []*github.RepositoryCommit{
		mockCommit("feat: breaking changes text capitalized to call it out strongly", "Bob", "bob", "38f1bd1091e162416bbcc653da5865b8f70e2c49"),
		mockCommit("feat: date on version is simplified", "Charlie", "charlie", "72372265d197605918b127c92eb75375c3715382"),
		mockCommit("feat: extract prerelease identifier action", "Alice", "alice", "0ba489f5f33d221061c149fed64166c26c6322ae"),
	}

	short, full, err := WriteChangelog("jakbytes", "version_actions", semver.MustParse("v0.0.0"), semver.MustParse("v0.1.0-src.0"), conventional.Commits{
		Feat: feat,
	}, false)
	require.Nil(t, err)

	var expectedShort = []string{
		"## [v0.1.0-src.0](https://github.com/jakbytes/version_actions/compare/v0.0.0...v0.1.0-src.0) (2024-02-07)",
		"### Features",
		"",
		"- ([`38f1bd1`](https://github.com/jakbytes/version_actions/commit/38f1bd1091e162416bbcc653da5865b8f70e2c49)) breaking changes text capitalized to call it out strongly",
		"- ([`7237226`](https://github.com/jakbytes/version_actions/commit/72372265d197605918b127c92eb75375c3715382)) date on version is simplified",
		"- ([`0ba489f`](https://github.com/jakbytes/version_actions/commit/0ba489f5f33d221061c149fed64166c26c6322ae)) extract prerelease identifier action",
		"",
	}

	for i, line := range expectedShort {
		if strings.HasPrefix(line, "## [") {
			prefix := strings.Split(line, ") (")[0]
			assert.True(t, strings.HasPrefix(short[i], prefix))
		} else {
			assert.Equal(t, line, short[i])
		}
	}

	var expectedLong = []string{
		"# Changelog",
		"",
		"## [v0.1.0-src.0](https://github.com/jakbytes/version_actions/compare/v0.0.0...v0.1.0-src.0) (2024-02-07)",
		"### Features",
		"",
		"- ([`38f1bd1`](https://github.com/jakbytes/version_actions/commit/38f1bd1091e162416bbcc653da5865b8f70e2c49)) breaking changes text capitalized to call it out strongly",
		"- ([`7237226`](https://github.com/jakbytes/version_actions/commit/72372265d197605918b127c92eb75375c3715382)) date on version is simplified",
		"- ([`0ba489f`](https://github.com/jakbytes/version_actions/commit/0ba489f5f33d221061c149fed64166c26c6322ae)) extract prerelease identifier action",
		"",
		"## [v0.0.0] Initial Version (2024-02-07)",
		"### Features",
		"",
		"- ([`c0d1dcd`](https://github.com/jakbytes/version_actions/commit/c0d1dcd0e3483390d8d7405569bcf3eadcce5710)) initial supported actions, version, sync, pull_request, extract_commit, download_release_asset",
		"",
		"### Fixes",
		"",
		"- ([`58bf05c`](https://github.com/jakbytes/version_actions/commit/58bf05caf571984ec6b2233ddb6f18a109a624ba)) type value needs to be output for further activity",
		"- ([`e1729a9`](https://github.com/jakbytes/version_actions/commit/e1729a947a61a321155939e72779334c88033b47)) action trigger should be set properly",
		"- ([`ba3d06f`](https://github.com/jakbytes/version_actions/commit/ba3d06fc58c65dc4fae5dd39c0d539207d906118)) hanging % needed to be removed from version action",
		"- ([`19bfb4d`](https://github.com/jakbytes/version_actions/commit/19bfb4db2aa5af63bead5067d2d3582e6b67fba2)) don't use best effort",
		"- ([`1a481d7`](https://github.com/jakbytes/version_actions/commit/1a481d72d0715ae6d7d88a9b434502513529c18c)) should be using v4 actions checkout",
		"- ([`1487ff3`](https://github.com/jakbytes/version_actions/commit/1487ff34f740541c9cb5aa3345aa14e6d1d93abc)) commits should be freeform to allow release and others",
		"- ([`68906c8`](https://github.com/jakbytes/version_actions/commit/68906c816d30d62c6f67c4a35b5e6003ccd74fbf)) download_release_asset shouldnt have quotes around the chmod val, version should not modify yml",
		"- ([`42328c0`](https://github.com/jakbytes/version_actions/commit/42328c0dc7d95b59e58c1373f678834420f8c329)) actions should reference version_action, not action",
		"- ([`8d24825`](https://github.com/jakbytes/version_actions/commit/8d24825ef39953f45c2fae275b420777c635ba5c)) a few more references to the old path were not adjusted",
		"- ([`db31802`](https://github.com/jakbytes/version_actions/commit/db31802dc409e7306ca2a4b17a8a1ba3e8332c05)) use the download_release_asset in pull_request, rename action.go to version_action.go",
		"",
		"### CI/CD",
		"",
		"- ([`a48f0ae`](https://github.com/jakbytes/version_actions/commit/a48f0aeac3a5c4ce3bed5af4e055bff7174bd99f)) fix reference to type",
		"- ([`3b55e7f`](https://github.com/jakbytes/version_actions/commit/3b55e7fbce860c789836006c2c1e93ab3a1554ce)) actions need to reference the correct path",
		"- ([`ed5f7a3`](https://github.com/jakbytes/version_actions/commit/ed5f7a398dd060d3a9769c344206c2b86dad2959)) remove debugging action",
	}

	require.Equal(t, len(expectedLong), len(full))

	for i, line := range expectedShort {
		if strings.HasPrefix(line, "## [") {
			prefix := strings.Split(line, ") (")[0]
			assert.True(t, strings.HasPrefix(short[i], prefix))
		} else {
			assert.Equal(t, line, short[i])
		}
	}

}
