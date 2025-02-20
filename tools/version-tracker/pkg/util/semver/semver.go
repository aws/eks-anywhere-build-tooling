package semver

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/constants"
)

type Version struct {
	Major, Minor, Patch       int64
	Prerelease, Buildmetadata string
}

func New(version string) (*Version, error) {
	semverRegexp := regexp.MustCompile(constants.SemverRegex)
	matches := semverRegexp.FindStringSubmatch(version)
	namedGroups := make(map[string]string, len(matches))
	groupNames := semverRegexp.SubexpNames()
	for i, value := range matches {
		name := groupNames[i]
		if name != "" {
			namedGroups[name] = value
		}
	}

	v := &Version{}
	var err error

	v.Major, err = strconv.ParseInt(namedGroups["major"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid major version in semver %s: %v", version, err)
	}
	v.Minor, err = strconv.ParseInt(namedGroups["minor"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid minor version in semver %s: %v", version, err)
	}
	v.Patch, err = strconv.ParseInt(namedGroups["patch"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid patch version in semver %s: %v", version, err)
	}

	v.Prerelease = namedGroups["prerelease"]
	v.Buildmetadata = namedGroups["buildmetadata"]

	return v, nil
}

func (v *Version) SameMajor(v2 *Version) bool {
	return v.Major == v2.Major
}

func (v *Version) SameMinor(v2 *Version) bool {
	return v.SameMajor(v2) && v.Minor == v2.Minor
}

func (v *Version) SamePatch(v2 *Version) bool {
	return v.SameMinor(v2) && v.Patch == v2.Patch
}

func (v *Version) SamePrerelease(v2 *Version) bool {
	return v.SamePatch(v2) && v.Prerelease == v2.Prerelease
}

func (v *Version) Equal(v2 *Version) bool {
	return v.SamePrerelease(v2) && v.Buildmetadata == v2.Buildmetadata
}

func (v *Version) GreaterThan(v2 *Version) bool {
	return v.Compare(v2) == 1
}

func (v *Version) Compare(v2 *Version) int {
	if c := compare(v.Major, v2.Major); c != 0 {
		return c
	}
	if c := compare(v.Minor, v2.Minor); c != 0 {
		return c
	}
	if c := compare(v.Patch, v2.Patch); c != 0 {
		return c
	}
	return 0
}

func compare(i, i2 int64) int {
	if i > i2 {
		return 1
	} else if i < i2 {
		return -1
	}
	return 0
}
