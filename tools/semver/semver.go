package semver

import (
	"github.com/Masterminds/semver/v3"
	"strconv"
	"strings"
)

type Version struct {
	*semver.Version
}

func (v *Version) PrereleaseIdentifier() string {
	return strings.Split(v.Prerelease(), ".")[0]
}

func (v *Version) PrereleaseVersionNumber() (int, error) {
	return strconv.Atoi(strings.Split(v.Prerelease(), ".")[1])
}

func (v *Version) IsPrerelease() bool {
	return v.Prerelease() != ""
}

func NewVersion(version string) (*Version, error) {
	v, err := semver.NewVersion(version)
	return &Version{v}, err
}

func MustParse(version string) *Version {
	v := semver.MustParse(version)
	return &Version{v}
}

func (v *Version) IncMajor() *Version {
	ver := v.Version.IncMajor()
	return &Version{&ver}
}

func (v *Version) IncMinor() *Version {
	ver := v.Version.IncMinor()
	return &Version{&ver}
}

func (v *Version) IncPatch() *Version {
	ver := v.Version.IncPatch()
	return &Version{&ver}
}

func (v *Version) AsPrereleaseVersion(prereleaseIdentifier string, versionNumber int) (*Version, error) {
	return NewVersion(v.String() + "-" + prereleaseIdentifier + "." + strconv.Itoa(versionNumber))
}

func (v *Version) GreaterThan(version *Version) bool {
	return v.Version.GreaterThan(version.Version)
}
func (v *Version) LessThan(version *Version) bool {
	return v.Version.LessThan(version.Version)
}
