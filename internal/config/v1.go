package config

import "github.com/AntonKosov/git-backups/internal/slice"

type v1 struct {
	Repositories struct {
		Generic []generic `yaml:"generic"`
		GitHub  []gitHub  `yaml:"github"`
	} `yaml:"repositories"`
}

type generic struct {
	Name       string   `yaml:"profile"`
	RootFolder string   `yaml:"root_folder"`
	Targets    []target `yaml:"targets"`
}

type target struct {
	URL    string `yaml:"url"`
	Folder string `yaml:"folder"`
}

type gitHub struct {
	Name        string   `yaml:"profile"`
	RootFolder  string   `yaml:"root_folder"`
	Affiliation string   `yaml:"affiliation"`
	Token       string   `yaml:"token"`
	Include     []string `yaml:"include"`
	Exclude     []string `yaml:"exclude"`
}

func (v v1) transform() Config {
	return Config{
		Repositories: Repositories{
			Generic: slice.Map(v.Repositories.Generic, func(g generic) GenericRepo {
				return GenericRepo{
					Name:       g.Name,
					RootFolder: g.RootFolder,
					Targets: slice.Map(g.Targets, func(t target) GenericTarget {
						return GenericTarget(t)
					}),
				}
			}),
			GitHub: slice.Map(v.Repositories.GitHub, func(g gitHub) GitHubRepo {
				return GitHubRepo(g)
			}),
		},
	}
}
