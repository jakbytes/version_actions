package changelog

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/google/go-github/v58/github"
	"github.com/jakbytes/version_actions/internal/utility"
	"github.com/jakbytes/version_actions/tools/conventional"
	"github.com/jakbytes/version_actions/tools/semver"
	"io/fs"
	"os"
	"strings"
	"time"
)

type Markdown []string

func (b Markdown) String() (s string) {
	var sb strings.Builder
	for i, line := range b {
		if i != len(b)-1 {
			sb.WriteString(line + "\n")
		} else { // last line
			sb.WriteString(line)
		}
	}
	return sb.String()
}

var Path = "CHANGELOG.md"

type Section struct {
	Title   string
	Commits []*github.RepositoryCommit
}

// GenerateNewChangelog generates a Markdown formatted changelog from the provided GitHub commits. It is intended to
// aggregate the changes from just the commits since the previous version.
func GenerateNewChangelog(org, repo string, previousVersion, version *semver.Version, commits conventional.Commits, disableVersionHeader bool) (body Markdown) {
	body = append(body, generateVersionHeader(org, repo, previousVersion, version, disableVersionHeader))

	sections := []Section{
		{"⚠ BREAKING CHANGES", commits.Breaking},
		{"Features", commits.Feat},
		{"Fixes", commits.Fix},
		{"Documentation", commits.Docs},
		{"Styles", commits.Style},
		{"Refactors", commits.Refactor},
		{"Performance", commits.Perf},
		{"Test", commits.Test},
		{"Build", commits.Build},
		{"CI/CD", commits.CI},
		{"Debugging", commits.Debug},
	}

	for _, section := range sections {
		if len(section.Commits) > 0 {
			body = append(body, fmt.Sprintf("### %s", section.Title), "")
			for _, commit := range section.Commits {
				body = append(body, formatCommit(org, repo, commit)...)
			}
			body = append(body, "")
		}
	}

	return
}

func generateVersionHeader(org, repo string, previousVersion, version *semver.Version, disableVersionHeader bool) string {
	currentDate := time.Now().UTC().Format("2006-01-02")

	if disableVersionHeader {
		return "## Changelog"
	} else if previousVersion != nil {
		// Header for the version with GitHub compare link
		return fmt.Sprintf("## [v%s](https://github.com/%s/%s/compare/v%s...v%s) (%s)", version, org, repo, previousVersion, version, currentDate)
	} else {
		return fmt.Sprintf("## [v%s] Initial Version (%s)", version, currentDate)
	}
}

func formatCommit(org, repo string, commit *github.RepositoryCommit) Markdown {
	// Extracting the first line of the commit message
	message := strings.TrimSpace(strings.SplitN(strings.TrimSpace(*commit.Commit.Message), ":", 2)[1])
	messageParts := strings.Split(message, "\n")
	// Extracting a short commit hash
	shortSHA := (*commit.SHA)[:7]

	m := Markdown{
		fmt.Sprintf("- ([`%s`](https://github.com/%s/%s/commit/%s)) %s", shortSHA, org, repo, *commit.SHA, messageParts[0]),
	}

	for _, line := range messageParts[1:] {
		m = append(m, fmt.Sprintf("  > %s", line))
	}
	return m
}

// Determines whether a line should be skipped.
func skipLine(line string, currentVersion, skipNextBreak, skipNextSpace *bool, versionHeading string) bool {
	if strings.HasPrefix(line, "# Changelog") {
		*skipNextBreak = true
		*skipNextSpace = true
		return true
	}

	if strings.HasPrefix(line, versionHeading) {
		*currentVersion = true
		return true
	} else if strings.HasPrefix(line, "## [") {
		*currentVersion = false
	}

	if strings.HasPrefix(line, "---") {
		if *skipNextBreak {
			*skipNextBreak = false
			*skipNextSpace = true
			return true
		}
	}

	if *skipNextSpace {
		*skipNextSpace = false
		return true
	}

	return *currentVersion
}

func updateChangelog(version *semver.Version, lines Markdown) (Markdown, error) {
	err := utility.Open(Path, func(file *os.File) (err error) {
		versionHeading := "## [v" + strings.Split(version.String(), "-")[0]
		currentVersion := false
		skipNextBreak := false
		skipNextSpace := false
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			skip := skipLine(line, &currentVersion, &skipNextBreak, &skipNextSpace, versionHeading)
			if skip {
				continue
			}

			lines = append(lines, line)
		}
		return nil
	})

	return lines, err
}

var UpdateChangelog = updateChangelog

func WriteChangelog(org, repo string, previousVersion, version *semver.Version, commits conventional.Commits, disableVersionHeader bool) (Markdown, Markdown, error) {
	changelog := GenerateNewChangelog(org, repo, previousVersion, version, commits, disableVersionHeader)
	lines := append(Markdown{"# Changelog", ""}, changelog...) // initialize lines with the header and version changelog
	_, err := os.Stat(Path)
	if !errors.Is(err, fs.ErrNotExist) { // CHANGELOG.md exists, update the file with the new version changelog and retain the rest of the file
		lines, err = UpdateChangelog(version, lines)
		if err != nil {
			return nil, nil, err
		}
	}
	return changelog, lines, WriteToFile(Path, lines)
}

func writeString(file *os.File, line string) error {
	_, err := file.WriteString(line + "\n")
	return err
}

var WriteString = writeString

func WriteToFile(path string, lines Markdown) error {
	return utility.Create(path, func(file *os.File) error {
		for _, line := range lines {
			err := WriteString(file, line)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
