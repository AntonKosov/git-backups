package config

import "github.com/AntonKosov/git-backups/internal/slice"

type v1 struct {
	Profiles struct {
		Generic []genericProfile `yaml:"generic"`
		GitHub  []gitHubProfile  `yaml:"github"`
	} `yaml:"profiles"`
}

type genericProfile struct {
	Name          string   `yaml:"profile"`
	RootFolder    string   `yaml:"root_folder"`
	PrivateSSHKey *string  `yaml:"private_ssh_key"`
	Targets       []target `yaml:"targets"`
}

type target struct {
	URL    string `yaml:"url"`
	Folder string `yaml:"folder"`
}

type gitHubProfile struct {
	Name          string   `yaml:"profile"`
	RootFolder    string   `yaml:"root_folder"`
	Affiliation   string   `yaml:"affiliation"`
	Token         string   `yaml:"token"`
	PrivateSSHKey *string  `yaml:"private_ssh_key"`
	Include       []string `yaml:"include"`
	Exclude       []string `yaml:"exclude"`
}

func (v v1) transform() Config {
	return Config{
		Profiles: Profiles{
			GenericProfiles: slice.Map(v.Profiles.Generic, func(g genericProfile) GenericProfile {
				return GenericProfile{
					Name:          g.Name,
					RootFolder:    g.RootFolder,
					PrivateSSHKey: g.PrivateSSHKey,
					Targets: slice.Map(g.Targets, func(t target) GenericTarget {
						return GenericTarget(t)
					}),
				}
			}),
			GitHubProfiles: slice.Map(v.Profiles.GitHub, func(g gitHubProfile) GitHubProfile {
				return GitHubProfile(g)
			}),
		},
	}
}
