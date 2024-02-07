package conventional

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateCommitMessage(t *testing.T) {
	// Table driven test with varying conventional commit message types
	tests := []struct {
		commitType    string
		commitMessage string
		expected      bool
	}{
		{"docs", "docs: add new documentation", true},
		{"docs", "docs(scope): add new documentation", true},

		// varying commit types
		{"feat", "docs!: add new documentation", false},
		{"feat", "docs: add new documentation", false},
		{"fix", "docs: add new documentation", false},

		// styles
		{"style", "style: code style", true},
		{"style", "style(scope): code style", true},

		// refactor
		{"refactor", "refactor: code refactor", true},
		{"refactor", "refactor(scope): code refactor", true},

		// perf
		{"perf", "perf: code perf", true},
		{"perf", "perf(scope): code perf", true},

		// test
		{"test", "test: code test", true},
		{"test", "test(scope): code test", true},

		// build
		{"build", "build: code build", true},
		{"build", "build(scope): code build", true},

		// ci
		{"ci", "ci: code ci", true},
		{"ci", "ci(scope): code ci", true},
	}

	for _, test := range tests {
		t.Run(test.commitType, func(t *testing.T) {
			require.Equal(t, test.expected, validateCommitMessage(test.commitType, test.commitMessage))
		})
	}
}
