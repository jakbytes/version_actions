package conventional

import (
	"errors"
	"fmt"
	"github.com/jakbytes/version_actions/tools/semver"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBumpVersion is a table-driven test for the BumpVersion function
func TestBumpVersion(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name                    string
		currentVersion          *semver.Version
		currentReleaseCandidate *semver.Version
		defaultBranch           string
		currentBranch           string
		increment               Increment
		prerelease              string
		want                    *semver.Version
		wantErr                 bool
		err                     error
	}{
		{"Initial Bump", nil, nil, "main", "main", Patch, "prerelease", semver.MustParse("0.0.0"), false, nil},
		{"No Increment", semver.MustParse("1.2.3"), nil, "main", "main", -1, "prerelease", semver.MustParse("1.2.3"), false, nil},
		{"Invalid Prerelease Version", semver.MustParse("1.2.3"), semver.MustParse("1.2.4-prerelease.G"), "main", "development", Patch, "prerelease", semver.MustParse("1.2.4-prerelease.1"), true, &strconv.NumError{
			Func: "Atoi",
			Num:  "G",
			Err:  errors.New("invalid syntax"),
		}},

		{"Prerelease Patch", semver.MustParse("1.2.3"), nil, "main", "development", Patch, "prerelease", semver.MustParse("1.2.4-prerelease.0"), false, nil},
		{"Prerelease Minor", semver.MustParse("1.2.3"), nil, "main", "development", Minor, "prerelease", semver.MustParse("1.3.0-prerelease.0"), false, nil},
		{"Prerelease Major", semver.MustParse("1.2.3"), nil, "main", "development", Major, "prerelease", semver.MustParse("2.0.0-prerelease.0"), false, nil},

		{"Prerelease Patch Increment", semver.MustParse("1.2.3"), semver.MustParse("1.2.4-prerelease.0"), "main", "development", Patch, "prerelease", semver.MustParse("1.2.4-prerelease.1"), false, nil},
		{"Prerelease Minor Increment", semver.MustParse("1.2.3"), semver.MustParse("1.3.0-prerelease.0"), "main", "development", Minor, "prerelease", semver.MustParse("1.3.0-prerelease.1"), false, nil},
		{"Prerelease Major Increment", semver.MustParse("1.2.3"), semver.MustParse("2.0.0-prerelease.0"), "main", "development", Major, "prerelease", semver.MustParse("2.0.0-prerelease.1"), false, nil},

		{"Prerelease Minor Increment from Patch", semver.MustParse("1.2.3"), semver.MustParse("1.2.4-prerelease.0"), "main", "development", Minor, "prerelease", semver.MustParse("1.3.0-prerelease.0"), false, nil},
		{"Prerelease Major Increment from Patch", semver.MustParse("1.2.3"), semver.MustParse("1.2.4-prerelease.0"), "main", "development", Major, "prerelease", semver.MustParse("2.0.0-prerelease.0"), false, nil},

		{"Prerelease Tag Change", semver.MustParse("1.2.3"), semver.MustParse("1.2.4-prerelease.0"), "main", "development", Patch, "alpha", semver.MustParse("1.2.4-alpha.0"), false, nil},
		{"Prerelease Tag Change with Increment", semver.MustParse("1.2.3"), semver.MustParse("1.2.4-prerelease.0"), "main", "development", Minor, "alpha", semver.MustParse("1.3.0-alpha.0"), false, nil},

		{"Default Patch Increment", semver.MustParse("1.2.3"), nil, "main", "main", Patch, "prerelease", semver.MustParse("1.2.4"), false, nil},
		{"Default Minor Increment", semver.MustParse("1.2.3"), nil, "main", "main", Minor, "prerelease", semver.MustParse("1.3.0"), false, nil},
		{"Default Major Increment", semver.MustParse("1.2.3"), nil, "main", "main", Major, "prerelease", semver.MustParse("2.0.0"), false, nil},
	}

	// Iterate over the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := IncVersion(VersionInfo{tc.currentVersion, tc.currentReleaseCandidate}, VersionConfig{tc.defaultBranch, tc.currentBranch, tc.prerelease}, tc.increment)
			if tc.wantErr {
				require.Error(t, err, "BumpVersion() should have returned an error")
				require.Equal(t, tc.err, err)
			} else {
				require.NoError(t, err, fmt.Sprintf("BumpVersion() should not have returned an error: %v", err))
				require.Equal(t, tc.want.String(), got.String(), "BumpVersion() did not return expected version")
			}
		})
	}
}
