package conventional

import (
	"github.com/jakbytes/version_actions/tools/semver"
)

// VersionConfig is a struct that contains the configuration for the versioning process
type VersionConfig struct {
	DefaultBranch        string // The default branch for the repository
	BaseBranch           string // The current branch version to increment
	PrereleaseIdentifier string // The prerelease tag to use for prerelease versions
}

// VersionInfo is a struct that contains the current version and the current release candidate version
type VersionInfo struct {
	CurrentVersion          *semver.Version
	CurrentReleaseCandidate *semver.Version
}

// IncVersion increments the version based on the current version, the current release candidate, and the configuration
func IncVersion(info VersionInfo, config VersionConfig, increment Increment) (*semver.Version, error) {

	var newVersion *semver.Version
	if info.CurrentVersion == nil {
		newVersion = semver.MustParse("0.0.0")
	} else {
		newVersion = incrementVersion(info.CurrentVersion, increment)
	}

	if config.BaseBranch != config.DefaultBranch {
		return incPrereleaseVersion(newVersion, info.CurrentReleaseCandidate, config.PrereleaseIdentifier)
	}

	return newVersion, nil
}

// incrementVersion increments the version based on the current version and the increment type
func incrementVersion(version *semver.Version, increment Increment) *semver.Version {
	switch increment {
	case Major:
		return version.IncMajor()
	case Minor:
		return version.IncMinor()
	case Patch:
		return version.IncPatch()
	default:
		// increment is likely -1 because no version increment was necessary
		return version
	}
}

func incPrereleaseVersion(newVersion *semver.Version, curReleaseCandidate *semver.Version, prereleaseIdentifier string) (*semver.Version, error) {
	if curReleaseCandidate == nil || (curReleaseCandidate.Major() != newVersion.Major() || curReleaseCandidate.Minor() != newVersion.Minor()) {
		// No release candidate exists or the current release candidate version is not the same as the new version
		return newVersion.AsPrereleaseVersion(prereleaseIdentifier, 0)
	} else {
		rcValue, err := curReleaseCandidate.PrereleaseVersionNumber()
		if err != nil {
			return nil, err
		}
		if curReleaseCandidate.PrereleaseIdentifier() == prereleaseIdentifier {
			return newVersion.AsPrereleaseVersion(prereleaseIdentifier, rcValue+1)
		} else {
			return newVersion.AsPrereleaseVersion(prereleaseIdentifier, 0)
		}
	}
}
