package conventional

import (
	"testing"
	"time"

	"github.com/google/go-github/v58/github"
	"github.com/stretchr/testify/assert"
)

func TestParseCommits(t *testing.T) {
	mockCommit := func(message string, date time.Time) *github.RepositoryCommit {
		return &github.RepositoryCommit{
			Commit: &github.Commit{
				Message: &message,
				Committer: &github.CommitAuthor{
					Date: &github.Timestamp{Time: date},
				},
			},
		}
	}

	// Test cases
	testCases := []struct {
		name           string
		commits        map[string]*github.RepositoryCommit
		expectBump     Increment
		expectBreaking int
		expectFeat     int
		expectFix      int
		expectDocs     int
		expectStyle    int
		expectRefactor int
		expectPerf     int
		expectTest     int
		expectBuild    int
		expectCI       int
	}{
		{
			name: "Commit with breaking change",
			commits: map[string]*github.RepositoryCommit{
				"1": mockCommit("feat: allow provided config object to extend other configs\n\nBREAKING CHANGE: `extends` key in config file is now used for extending other config files", time.Now()),
			},
			expectBump:     Major,
			expectBreaking: 1,
		},
		{
			name: "Commit with ! for breaking change",
			commits: map[string]*github.RepositoryCommit{
				"2": mockCommit("feat!: send an email to the customer when a product is shipped", time.Now()),
			},
			expectBump:     Major,
			expectBreaking: 1,
		},
		{
			name: "Commit with scope and ! for breaking change",
			commits: map[string]*github.RepositoryCommit{
				"3": mockCommit("feat(api)!: send an email to the customer when a product is shipped", time.Now()),
			},
			expectBump:     Major,
			expectBreaking: 1,
		},
		{
			name: "Commit with both ! and BREAKING CHANGE footer",
			commits: map[string]*github.RepositoryCommit{
				"4": mockCommit("chore!: drop support for Node 6\n\nBREAKING CHANGE: use JavaScript features not available in Node 6.", time.Now()),
			},
			expectBump:     Major,
			expectBreaking: 1,
		},
		{
			name: "Commit with no body",
			commits: map[string]*github.RepositoryCommit{
				"5": mockCommit("docs: correct spelling of CHANGELOG", time.Now()),
			},
			expectBump:     -1,
			expectBreaking: 0,
			expectDocs:     1,
		},
		{
			name: "Commit with scope",
			commits: map[string]*github.RepositoryCommit{
				"6": mockCommit("feat(lang): add Polish language", time.Now()),
			},
			expectBump: Minor,
			expectFeat: 1,
		},
		{
			name: "Commit with multi-paragraph body and multiple footers",
			commits: map[string]*github.RepositoryCommit{
				"7": mockCommit("fix: prevent racing of requests\n\nIntroduce a request id and a reference to latest request. Dismiss\nincoming responses other than from latest request.\n\nRemove timeouts which were used to mitigate the racing issue but are\nobsolete now.\n\nReviewed-by: Z\nRefs: #123", time.Now()),
			},
			expectBump: Patch,
			expectFix:  1,
		},
		{
			name: "Commit with malformed message",
			commits: map[string]*github.RepositoryCommit{
				"8": mockCommit("This is a malformed commit message without conventional structure.", time.Now()),
			},
			expectBump: -1,
			// Since the parsing is expected to fail, there should be no categorization of this commit.
			expectBreaking: 0,
			expectFeat:     0,
			expectFix:      0,
		},
		{
			name: "Mixed major, minor, and patch commits",
			commits: map[string]*github.RepositoryCommit{
				"2": mockCommit("feat: implement user profiles", time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)),  // Minor
				"6": mockCommit("fix: fix login bug", time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC)),             // Patch
				"3": mockCommit("feat!: add new logging feature", time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)), // Major due to breaking change
				"4": mockCommit("feat: implement user logout", time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)),    // Minor
				"5": mockCommit("feat: implement user records", time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),   // Minor
				"1": mockCommit("feat: implement user login", time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)),     // Minor
			},
			expectBump:     Major,
			expectBreaking: 1,
			expectFeat:     4,
			expectFix:      1,
		},
		{
			name: "all commit types",
			commits: map[string]*github.RepositoryCommit{
				"2": mockCommit("feat: implement user profiles", time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)),  // Minor
				"6": mockCommit("fix: fix login bug", time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC)),             // Patch
				"3": mockCommit("feat!: add new logging feature", time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)), // Major due to breaking change
				"4": mockCommit("feat: implement user logout", time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)),    // Minor
				"5": mockCommit("feat: implement user records", time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),   // Minor
				"1": mockCommit("feat: implement user login", time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)),     // Minor
				// all other types
				"7":  mockCommit("docs: correct spelling of CHANGELOG", time.Now()),
				"8":  mockCommit("style: add missing semicolons", time.Now()),
				"9":  mockCommit("refactor: share logic between 4d3d3d3 and 2b2b2b2", time.Now()),
				"10": mockCommit("perf: remove O(n) algorithm", time.Now()),
				"11": mockCommit("test: add missing unit tests", time.Now()),
				"12": mockCommit("build: add build script", time.Now()),
				"13": mockCommit("ci: add CI script", time.Now()),
			},
			expectBump:     Major,
			expectBreaking: 1,
			expectFeat:     4,
			expectFix:      1,
			expectDocs:     1,
			expectStyle:    1,
			expectRefactor: 1,
			expectPerf:     1,
			expectTest:     1,
			expectBuild:    1,
			expectCI:       1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parsed := ParseCommits(tc.commits)
			assert.Equal(t, tc.expectBump, parsed.Increment())
			assert.Len(t, parsed.Breaking, tc.expectBreaking)
			assert.Len(t, parsed.Feat, tc.expectFeat)
			assert.Len(t, parsed.Fix, tc.expectFix)
			assert.Len(t, parsed.Docs, tc.expectDocs)
			assert.Len(t, parsed.Style, tc.expectStyle)
			assert.Len(t, parsed.Refactor, tc.expectRefactor)
			assert.Len(t, parsed.Perf, tc.expectPerf)
			assert.Len(t, parsed.Test, tc.expectTest)
			assert.Len(t, parsed.Build, tc.expectBuild)
			assert.Len(t, parsed.CI, tc.expectCI)

			// Assuming all parsed are merged into a single slice for sorting validation
			for _, commits := range [][]*github.RepositoryCommit{parsed.Breaking, parsed.Feat, parsed.Fix} {
				for i := 0; i < len(commits)-1; i++ {
					assert.True(t, commits[i].Commit.Committer.Date.After(*commits[i+1].Commit.Committer.Date.GetTime()))
				}
			}
		})
	}
}

func TestInsertDescendingOrder(t *testing.T) {
	// Define the comparison function for descending order
	descending := func(i, j int) bool {
		return i > j
	}

	// Initialize an unsorted slice
	x := []int{1, 5, 3, 2, 4}
	var o []int
	// Insert items in unsorted order
	for _, i := range x {
		o = insert(o, i, descending)
	}

	// Expected result is [5, 4, 3, 2, 1]
	expected := []int{5, 4, 3, 2, 1}

	// Use assert.Equal to check if the result matches the expected result
	assert.Equal(t, expected, o)
}
