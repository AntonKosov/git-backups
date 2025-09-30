package config

type Config struct {
	Profiles Profiles
}

type Profiles struct {
	GenericProfiles []GenericProfile
	GitHubProfiles  []GitHubProfile
}

type GenericProfile struct {
	Name          string
	RootFolder    string
	PrivateSSHKey *string
	Targets       []GenericTarget
}

type GenericTarget struct {
	URL    string
	Folder string
}

type GitHubProfile struct {
	Name          string
	RootFolder    string
	Affiliation   string
	Token         string
	PrivateSSHKey *string
	Include       []string
	Exclude       []string
}
