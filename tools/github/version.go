package github

import (
	"github.com/google/go-github/v58/github"
	"github.com/rs/zerolog/log"
	"version_actions/tools/semver"
)

// Version is a struct that contains the semantic version and the GitHub RepositoryTag associated with it.
type Version struct {
	*semver.Version
	*github.RepositoryTag
}

// Versions maintains a slice of Version structs and the latest version in the slice
type Versions struct {
	inner  []*Version
	latest *Version // latest version
}

// RepositoryVersions contains the stable and prerelease versions for a repository
type RepositoryVersions struct {
	release    *Versions            // stable versions
	prerelease map[string]*Versions // prerelease versions, keyed by prerelease name
}

// Tags returns the list of tags in the repository. The tags are cached in the repository struct, so subsequent calls
// to Tags will not make additional network requests.
func (r *Repository) Tags() (tags []*github.RepositoryTag, err error) {
	if r.tags == nil {
		r.tags, _, err = r.ListTags(r.Ctx, r.RepositoryMetadata.Owner, r.RepositoryMetadata.Name, nil)
	}
	return r.tags, err
}

// Versions returns the stable and prerelease versions for the repository. The tags and versions are cached in the
// repository struct, so subsequent calls to Versions will not make additional network requests. Tags that are not valid
// semantic versions are ignored.
func (r *Repository) Versions() (*RepositoryVersions, error) {
	if r.versions == nil {
		r.versions = &RepositoryVersions{&Versions{}, make(map[string]*Versions)}
	}
	if r.versions.release.inner == nil {
		tags, err := r.Tags()
		if err != nil {
			return &RepositoryVersions{}, err
		}

		for _, tag := range tags {
			r.parseTag(tag)
		}
	}
	return r.versions, nil
}

// parseTag parses the tag as a semantic version, if the tag is not a valid semantic version, it is ignored.
func (r *Repository) parseTag(tag *github.RepositoryTag) {
	version, err := semver.NewVersion(*tag.Name)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to parse version: %s", version)
		return
	}

	if version.IsPrerelease() {
		r.parsePrereleaseVersion(tag, version)
	} else {
		parseVersion(r.versions.release, tag, version)
	}
}

// parsePrereleaseVersion parses the prerelease version and adds it to the prerelease versions slice. If the prerelease
// version is greater than the latest prerelease version, it is set as the latest prerelease version.
func (r *Repository) parsePrereleaseVersion(tag *github.RepositoryTag, version *semver.Version) {
	pid := version.PrereleaseIdentifier()
	if _, ok := r.versions.prerelease[pid]; !ok {
		r.versions.prerelease[pid] = &Versions{}
	}
	parseVersion(r.versions.prerelease[version.PrereleaseIdentifier()], tag, version)
}

// parseVersion parses the version and adds it to the versions slice. If the version is greater than the latest version,
// it is set as the latest version.
func parseVersion(versions *Versions, tag *github.RepositoryTag, version *semver.Version) {
	versions.inner = append(versions.inner, &Version{
		Version:       version,
		RepositoryTag: tag,
	})

	if versions.latest == nil || version.GreaterThan(versions.latest.Version) {
		versions.latest = &Version{
			Version:       version,
			RepositoryTag: tag,
		}
	}
}

// LatestVersion returns the Version with the highest release tag in the repository. If there are no release tags, an
// error is returned. The tags and latest version are cached in the repository struct, so subsequent calls to
// LatestVersion will not make additional network requests. Invalid semantic version tags are ignored.
func (r *Repository) LatestVersion() (*Version, error) {
	versions, err := r.Versions()
	if err != nil {
		return &Version{}, err
	}

	if versions.release.latest == nil {
		return nil, NoReleaseVersionFound{}
	}

	return versions.release.latest, nil
}

// PreviousVersion returns the Version immediately preceding the Version with the highest release tag in the repository.
// If there are no release tags, an error is returned. The tags and latest version are cached in the repository struct,
//
//	so subsequent calls to PreviousVersion will not make additional network requests. Invalid semantic version tags are
//	ignored.
func (r *Repository) PreviousVersion() (*Version, error) {
	versions, err := r.Versions()
	if err != nil {
		return &Version{}, err
	}

	if versions.release.latest == nil {
		return &Version{}, NoReleaseVersionFound{}
	}
	var previous *Version
	for _, version := range versions.release.inner {
		if version.String() != versions.release.latest.String() && (previous == nil || version.GreaterThan(previous.Version)) {
			previous = version
		}
	}
	return previous, nil
}

// LatestPrereleaseVersion returns the Version with the highest prerelease tag in the repository with the given
// prerelease identifier. If there are no prerelease tags, an error is returned. The tags and latest version are cached
// in the repository struct, so subsequent calls to LatestPrereleaseVersion will not make additional network requests.
func (r *Repository) LatestPrereleaseVersion(prereleaseIdentifier string) (*Version, error) {
	versions, err := r.Versions()
	if err != nil {
		return nil, err
	}

	if versions.prerelease[prereleaseIdentifier] == nil {
		return &Version{}, NoPrereleaseVersionFound{}
	}

	return versions.prerelease[prereleaseIdentifier].latest, nil
}
