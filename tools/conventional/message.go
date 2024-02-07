package conventional

import (
	"fmt"
	"github.com/google/go-github/v58/github"
	"github.com/leodido/go-conventionalcommits"
	"regexp"
	"sync"
)

var (
	// Mutex for concurrent access to the regex map
	mutex sync.Mutex

	// Map to store compiled regex patterns
	regexMap = make(map[string]*regexp.Regexp)
)

type Message struct {
	*conventionalcommits.ConventionalCommit
	*github.RepositoryCommit
}

// getCompiledRegex returns a compiled regex for a given commit type.
// It compiles the regex if it's not already in the map.
func getCompiledRegex(commitType string) *regexp.Regexp {
	mutex.Lock()
	defer mutex.Unlock()

	// Check if the regex is already compiled
	if re, exists := regexMap[commitType]; exists {
		return re
	}

	// Compile the regex
	pattern := fmt.Sprintf("^%s(\\(.*\\))?:\\s.*$", regexp.QuoteMeta(commitType))
	re := regexp.MustCompile(pattern)

	// Store the compiled regex in the map
	regexMap[commitType] = re
	return re
}

// validateCommitMessage checks if a commit message conforms to the conventional commit format.
func validateCommitMessage(commitType, commitMessage string) bool {
	re := getCompiledRegex(commitType)
	return re.MatchString(commitMessage)
}

func (m *Message) IsDocs() bool {
	return validateCommitMessage("docs", *m.Commit.Message)
}

func (m *Message) IsStyle() bool {
	return validateCommitMessage("style", *m.Commit.Message)
}

func (m *Message) IsRefactor() bool {
	return validateCommitMessage("refactor", *m.Commit.Message)
}

func (m *Message) IsPerf() bool {
	return validateCommitMessage("perf", *m.Commit.Message)
}

func (m *Message) IsTest() bool {
	return validateCommitMessage("test", *m.Commit.Message)
}

func (m *Message) IsBuild() bool {
	return validateCommitMessage("build", *m.Commit.Message)
}

func (m *Message) IsCI() bool {
	return validateCommitMessage("ci", *m.Commit.Message)
}
