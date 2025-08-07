package config

type V1 struct {
	Repositories struct {
		Generic []struct {
			Name       string `yaml:"category"`
			RootFolder string `yaml:"root_folder"`
			Targets    []struct {
				URL    string `yaml:"url"`
				Folder string `yaml:"folder"`
			} `yaml:"targets"`
		} `yaml:"generic"`

		GitHub []struct {
			Name       string   `yaml:"category"`
			RootFolder string   `yaml:"root_folder"`
			Token      string   `yaml:"token"`
			Include    []string `yaml:"include"`
			Exclude    []string `yaml:"exclude"`
		} `yaml:"github"`
	} `yaml:"repositories"`
}
