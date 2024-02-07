package conventional

import (
	"github.com/google/go-github/v58/github"
	"github.com/jakbytes/version_actions/internal/logger"
	"github.com/leodido/go-conventionalcommits"
	"github.com/leodido/go-conventionalcommits/parser"
	"github.com/rs/zerolog/log"
	"sort"
)

// Increment is a type that represents the type of increment to make to the version.
type Increment int

const (
	Major Increment = iota
	Minor
	Patch
)

// Parser is a struct that contains the parser for conventional commit messages.
type Parser struct {
	conventionalcommits.Machine
}

// ParseCommit parses the commit message and returns the conventional commit message.
func (p *Parser) ParseCommit(commit *github.RepositoryCommit) (out *conventionalcommits.ConventionalCommit) {
	message, err := p.Parse([]byte(*commit.Commit.Message))
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to parse commit message: %s", *commit.Commit.Message)
	}
	out, _ = message.(*conventionalcommits.ConventionalCommit)
	return
}

// Commits is a struct that contains the commits for breaking changes, features, and fixes based on the conventional
// commit types.
type Commits struct {
	Breaking []*github.RepositoryCommit
	Feat     []*github.RepositoryCommit
	Fix      []*github.RepositoryCommit
	Docs     []*github.RepositoryCommit
	Style    []*github.RepositoryCommit
	Refactor []*github.RepositoryCommit
	Perf     []*github.RepositoryCommit
	Test     []*github.RepositoryCommit
	Build    []*github.RepositoryCommit
	CI       []*github.RepositoryCommit
	Debug    []*github.RepositoryCommit
}

// Increment returns the increment type based on the collection of commits.
// If there are any breaking changes, the increment type is Major. If there are any features, the increment type is
// Minor. If there are any fixes, the increment type is Patch. Otherwise, the increment type is -1, indicating no
// increment is necessary.
func (c *Commits) Increment() Increment {
	if len(c.Breaking) > 0 {
		return Major
	}
	if len(c.Feat) > 0 {
		return Minor
	}
	if len(c.Fix) > 0 {
		return Patch
	}
	return -1
}

// ParseCommits parses the commits and returns the version increment and associated commits for breaking changes, features,
// and fixes.
//
// The parser is configured to use the best effort mode. The best effort mode will make the parser return what it found
// until the point it errored out, if it found (at least) a valid type and a valid description. However, if the parser
// does not find a valid type or a valid description, it will not account for the commit.
//
// See: https://github.com/leodido/go-conventionalcommits?tab=readme-ov-file#best-effort
//
// Parameters:
//   - commits: The commits to parse.
//
// Returns:
//   - parsed (Commits): The parsed commits.
func ParseCommits(commits map[string]*github.RepositoryCommit) (parsed Commits) {
	log.Logger = logger.Base()
	cparser := Parser{parser.NewMachine(
		conventionalcommits.WithTypes(conventionalcommits.TypesFreeForm),
	)}
	for _, commit := range commits {
		message := Message{
			cparser.ParseCommit(commit),
			commit,
		}
		if message.ConventionalCommit != nil {
			if message.IsBreakingChange() {
				parsed.Breaking = insert(parsed.Breaking, commit, less)
			} else if message.IsFeat() {
				parsed.Feat = insert(parsed.Feat, commit, less)
			} else if message.IsFix() {
				parsed.Fix = insert(parsed.Fix, commit, less)
			} else if message.IsDocs() {
				parsed.Docs = insert(parsed.Docs, commit, less)
			} else if message.IsStyle() {
				parsed.Style = insert(parsed.Style, commit, less)
			} else if message.IsRefactor() {
				parsed.Refactor = insert(parsed.Refactor, commit, less)
			} else if message.IsPerf() {
				parsed.Perf = insert(parsed.Perf, commit, less)
			} else if message.IsTest() {
				parsed.Test = insert(parsed.Test, commit, less)
			} else if message.IsBuild() {
				parsed.Build = insert(parsed.Build, commit, less)
			} else if message.IsCI() {
				parsed.CI = insert(parsed.CI, commit, less)
			} else if message.IsDebug() {
				parsed.Debug = insert(parsed.Debug, commit, less)
			}

		}
	}
	return parsed
}

// less function returns true if the commit date of i is less than j, false otherwise.
func less(i, j *github.RepositoryCommit) bool {
	return i.Commit.Committer.Date.After(*j.Commit.Committer.Date.GetTime())
}

// insert inserts the value into the slice x while maintaining the order of the slice
func insert[T any](x []T, value T, less func(i, j T) bool) []T {
	insertionIndex := sort.Search(len(x), func(i int) bool { return less(value, x[i]) })
	newX := append([]T(nil), x[:insertionIndex]...) // Copy first part
	newX = append(newX, value)                      // Append the new value
	newX = append(newX, x[insertionIndex:]...)      // Append the rest
	return newX
}
