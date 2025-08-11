package config

type Config struct {
	Repositories Repositories
}

type Repositories struct {
	Generic []GenericRepo
	GitHub  []GitHubRepo
}

type GenericRepo struct {
	Name       string
	RootFolder string
	Targets    []GenericTarget
}

type GenericTarget struct {
	URL    string
	Folder string
}

type GitHubRepo struct {
	Name       string
	RootFolder string
	Token      string
	Include    []string
	Exclude    []string
}
