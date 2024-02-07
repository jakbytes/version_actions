package composite

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"regexp"
	"version_actions/internal/utility"
	"version_actions/tools/changelog"
	"version_actions/tools/conventional"
	"version_actions/tools/github"
	"version_actions/tools/semver"
)

type Handler struct {
	*github.Client
	Owner                string
	Name                 string
	Head                 string
	hb                   *github.Branch
	Base                 string
	bb                   *github.Branch
	PrereleaseIdentifier string
	ReleaseBranch        string
	Latest               *github.Version
	LatestPrerelease     *github.Version
	CommitFiles          []string

	commits         *conventional.Commits
	title           string
	body            changelog.Markdown
	latestChangelog changelog.Markdown
	fullChangelog   changelog.Markdown

	inner     error
	promotion bool

	Trigger string
}

func (h *Handler) Wrapper(f func() error) {
	if h.inner != nil {
		h.inner = f()
	}
}

const releaseBranchPrefix = "release--branch--"

func (h *Handler) gatherVersions() {
	var err error
	h.Latest, err = h.Repository().LatestVersion()
	if err != nil {
		if errors.Is(err, github.NoReleaseVersionFound{}) {
			log.Warn().Err(err).Msg("No release version found, semver-action will behave as if the next version should be v0.0.0")
		} else {
			panic(err)
		}
	}

	if h.PrereleaseIdentifier != "" {
		h.LatestPrerelease, err = h.Repository().LatestPrereleaseVersion(h.PrereleaseIdentifier)
		if err != nil && !errors.Is(err, github.NoPrereleaseVersionFound{}) {
			panic(err)
		}
	}
}

func (h *Handler) Commits() (commits *conventional.Commits) {
	if h.commits == nil {
		var sha *string
		if h.Latest != nil {
			sha = h.Latest.Commit.SHA
		}
		raw, err := h.head().GetCommitsSinceCommit(sha)
		if err != nil {
			panic(err)
		}
		c := conventional.ParseCommits(raw)
		h.commits = &c
	}
	return h.commits
}

func (h *Handler) base() *github.Branch {
	if h.bb == nil {
		var err error
		h.bb, err = h.Repository().Branch(h.Base)
		if err != nil {
			panic(err)
		}
	}
	return h.bb
}

func (h *Handler) head() *github.Branch {
	if h.hb == nil {
		branchName := releaseBranchPrefix + h.Head
		if h.Head != h.Base { // release branch generated off the base branch
			branchName = releaseBranchPrefix + h.Base
			h.promotion = true
		}
		h.setBranch(branchName)
	}
	return h.hb
}

// PullRequest creates a release--branch--{branchName} pull request for branchName
// PR Details:
// - title: "release({branchName}): {nextVersion}"
// - base: {branchName}
// - head: release--branch--{branchName}
// - body:
//   - if prerelease: ":robot: I have created a release candidate *beep* *boop*"
//   - else: ":robot: I have created a release *beep* *boop*"
func (h *Handler) PullRequest() error {
	log.Debug().Msgf("owner: %s", h.Owner)
	log.Debug().Msgf("name: %s", h.Name)
	log.Debug().Msgf("head: %s", h.Head)
	log.Debug().Msgf("base: %s", h.Base)
	log.Debug().Msgf("prereleaseIdentifier: %s", h.PrereleaseIdentifier)
	log.Debug().Msgf("releaseBranch: %s", h.ReleaseBranch)
	log.Debug().Msgf("base(): %s", h.base().Name)
	log.Debug().Msgf("head(): %s", h.head().Name)
	h.gatherVersions()
	h.gatherChangelog()
	h.composePullRequest()

	if h.promotion || h.Trigger != "release" {
		h.commitChangelog()
		h.setPullRequest()
	}

	err := changelog.WriteToFile("release.txt", h.latestChangelog)
	if err != nil {
		return err
	}

	return h.inner
}

func (h *Handler) setPullRequest() {
	err := h.SetPullRequest(h.head().Name, h.base().Name, h.title, false, func(_ *string) (changelog.Markdown, error) {
		return h.body, nil
	})
	if err != nil {
		panic(err)
	}
}

func (h *Handler) setBranch(name string) {
	head, err := h.Repository().Branch(h.Head)
	if err != nil {
		panic(err)
	}

	branch, err := h.Repository().Branch(name)
	log.Debug().Msgf("branch: %v, %v", branch, err)
	if errors.Is(err, github.BranchNotFound{Name: name}) {
		h.hb, err = h.Repository().CreateBranch(name, head.Commit.SHA)
	} else if err == nil {
		h.hb = branch
		err = h.hb.Reset(head.Commit.SHA)
	}
	if err != nil {
		panic(err)
	}
}

func (h *Handler) updateActionYaml(files []github.File) []github.File {
	filesToUpdate := []string{"pull_request/action.yml", "version/action.yml"}
	for _, file := range filesToUpdate {
		log.Debug().Msgf("Updating version in %s", file)
		value, err := updateVersionInYaml(file, "v"+h.NextVersion().String())
		if err != nil {
			panic(err)
		}
		files = append(files, github.File{Path: file, Content: value})
	}

	return files
}

// Function to update the version in the Yaml files
func updateVersionInYaml(filePath, newVersion string) (string, error) {
	// Read the content of the YML file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Convert the content to a string
	text := string(content)

	// Define a regex pattern to find the version string following `/tags/`
	pattern := `(/tags/)[^\s"]+`
	re := regexp.MustCompile(pattern)

	// Replace the found pattern with the new version
	updatedText := re.ReplaceAllString(text, "${1}"+newVersion)

	return updatedText, nil
}

func (h *Handler) updateAdditionalFiles(files []github.File) []github.File {
	for _, path := range h.CommitFiles {
		err := utility.Open(path, func(file *os.File) error {
			content, err := io.ReadAll(file)
			files = append(files, github.File{Path: path, Content: string(content)})
			return err
		})
		if err != nil {
			panic(err)
		}
	}
	return files
}

func (h *Handler) commitChangelog() {
	log.Info().Msg("Committing changelog")
	files := []github.File{{Path: "CHANGELOG.md", Content: h.fullChangelog.String()}}
	if h.Owner == "jakbytes" && h.Name == "version_actions" {
		files = h.updateActionYaml(files)
	}

	files = h.updateAdditionalFiles(files)

	newTreeSHA, parentCommitSHA, err := h.head().AddFiles(files)
	if err != nil {
		panic(err)
	}

	err = h.head().CommitChanges(newTreeSHA, parentCommitSHA, h.title)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) gatherChangelog() {
	var err error
	h.latestChangelog, h.fullChangelog, err = changelog.WriteChangelog(h.Owner, h.Name, h.VersionInfo().CurrentVersion, h.NextVersion(), *h.Commits(), false)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) composePullRequest() {
	h.title = fmt.Sprintf("release(%s): v%s", h.Base, h.NextVersion().String())
	header := ":robot: I have created a release candidate *beep* *boop*"
	if h.Base == h.ReleaseBranch { // if the release branch is the target, we're promoting a release candidate to a release
		header = ":robot: I have created a release *beep* *boop*"
	}

	h.body = append([]string{header,
		"",
		"---",
		"",
	}, h.latestChangelog...)
	// footer
	h.body = append(h.body, "---", "",
		"This release was composed by [version_actions](https://github.com/jakbytes/version_actions)")
}

func (h *Handler) VersionInfo() (info conventional.VersionInfo) {
	if h.Latest != nil {
		info.CurrentVersion = h.Latest.Version
	}
	if h.LatestPrerelease != nil {
		info.CurrentReleaseCandidate = h.LatestPrerelease.Version
	}
	return
}

func (h *Handler) NextVersion() *semver.Version {
	version, err := conventional.IncVersion(h.VersionInfo(),
		conventional.VersionConfig{
			DefaultBranch:        h.ReleaseBranch,
			BaseBranch:           h.Base,
			PrereleaseIdentifier: h.PrereleaseIdentifier,
		}, h.Commits().Increment())
	if err != nil {
		panic(err)
	}
	return version
}
