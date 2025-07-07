package buildinfo

var (
	// The git commit that was compiled.
	//It will be filled in by the compiler.
	gitCommit string

	// The semver git tag associated with this release.
	// It will be filled in by the compiler.
	version string
)

func GetCommit() string {
	return gitCommit
}

func GetVersion() string {
	return version
}
