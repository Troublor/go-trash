package storage

import "strconv"

const (
	VersionMajor = 1          // Major version component of the current release		VersionMajor = 1        // Major version component of the current release
	VersionMinor = 3          // Minor version component of the current release		VersionMinor = 9        // Minor version component of the current release
	VersionPatch = 1          // Patch version component of the current release		VersionPatch = 2        // Patch version component of the current release
	VersionMeta  = "unstable" // Version metadata to append to the version string		VersionMeta  = "stable" // Version metadata to append to the version string
)

func Version() string {
	return strconv.Itoa(VersionMajor) + "." + strconv.Itoa(VersionMinor) + "." + strconv.Itoa(VersionPatch) + "-" + VersionMeta
}
