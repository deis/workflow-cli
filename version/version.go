package version

// Version identifies this Deis product revision.
const Version = "2.1.0-dev"

// BuildVersion is the git revision of the build.
// Note: This value is overwritten by the linker during build
var BuildVersion = ""
