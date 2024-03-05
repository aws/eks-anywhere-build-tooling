package types

// DisplayOptions represents the options that can be passed to the `display` command.
type DisplayOptions struct {
	ProjectName        string
	PrintLatestVersion bool
}

// UpgradeOptions represents the options that can be passed to the `upgrade` command.
type UpgradeOptions struct {
	ProjectName string
	DryRun      bool
}

// ProjectsList represents the top-level projects list in the upstream projects tracker file.
type ProjectsList struct {
	Projects []Project `yaml:"projects"`
}

// Project represents a list of repositories that share the same owner or organization.
type Project struct {
	Org   string `yaml:"org"`
	Repos []Repo `yaml:"repos"`
}

// Repo represents a GitHub repository with a list of versions being built.
type Repo struct {
	Name     string    `yaml:"name"`
	Versions []Version `yaml:"versions"`
}

// Version represents a Git tag or commit for a repository and the Go version corresponding to the revision.
type Version struct {
	Tag       string `yaml:"tag,omitempty"`
	Commit    string `yaml:"commit,omitempty"`
	GoVersion string `yaml:"go_version,omitempty"`
}

// ProjectVersionInfo represents the current and latest revision for a project.
type ProjectVersionInfo struct {
	Org            string
	Repo           string
	CurrentVersion string
	LatestVersion  string
}

// ReleaseTarball represents the GitHub release asset name, binary name and related settings to get the
// Go version for a particular project.
type ReleaseTarball struct {
	AssetName                string
	BinaryName               string
	Extract                  bool
	OverrideAssetURL         string
	TrimLeadingVersionPrefix bool
}

// GoVersionSourceOfTruth represents the source of truth file and search string to get the Go version
// for a particular project.
type GoVersionSourceOfTruth struct {
	SourceOfTruthFile     string
	GoVersionSearchString string
}

type ImageMetadata struct {
	Tag         string `yaml:"tag,omitempty"`
	ImageDigest string `yaml:"imageDigest,omitempty"`
}

type EKSDistroRelease struct {
	Branch      string `json:"branch"`
	KubeVersion string `json:"kubeVersion"`
	Number      int    `json:"number"`
	Dev         *bool  `json:"dev,omitempty"`
}

type EKSDistroLatestReleases struct {
	Releases []EKSDistroRelease `json:"releases"`
	Latest   string             `json:"latest"`
}
