package config_test

import (
	"github.com/AntonKosov/git-backups/internal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config tests", func() {
	const configPath = "../../test/data/config.yaml"
	const invalidConfigPath = "../../test/data/invalid_config.yaml"

	var (
		configFile string
		conf       config.Config
		err        error
	)

	BeforeEach(func() {
		configFile = configPath
	})

	JustBeforeEach(func() {
		conf, err = config.ReadConfig(configFile)
	})

	It("does not return an error", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	It("parses config correctly", func() {
		Expect(conf).To(Equal(config.Config{
			Repositories: config.Repositories{
				Generic: []config.GenericRepo{
					{
						Name:       "profile name",
						RootFolder: "/home/user/git_backup/folder_name",
						Targets: []config.GenericTarget{
							{
								URL:    "https://github.com/Username1/repo_name_1.git",
								Folder: "repo_folder_name_1",
							},
							{
								URL:    "https://github.com/Username2/repo_name_2.git",
								Folder: "repo_folder_name_2",
							},
						},
					},
					{
						Name:       "profile name 2",
						RootFolder: "/home/user/git_backup/folder_name_2",
						Targets: []config.GenericTarget{
							{
								URL:    "https://github.com/Username3/repo_name_3.git",
								Folder: "repo_folder_name_3",
							},
							{
								URL:    "https://github.com/Username4/repo_name_4.git",
								Folder: "repo_folder_name_4",
							},
						},
					},
				},
				GitHub: []config.GitHubRepo{
					{
						Name:        "profile name 3",
						RootFolder:  "/home/user/git_backup/folder_name_3",
						Affiliation: "owner,collaborator,organization_member",
						Token:       "GH_XXX",
						Include: []string{
							"repo_name_1",
							"repo_name_2",
						},
						Exclude: []string{
							"repo_name_3",
						},
					},
					{
						Name:        "profile name 4",
						RootFolder:  "/home/user/git_backup/folder_name_4",
						Affiliation: "owner",
						Token:       "GH2_XXX",
						Include: []string{
							"repo_name_4",
							"repo_name_5",
						},
						Exclude: []string{
							"repo_name_6",
						},
					},
				},
			},
		}))
	})

	When("missing a config file", func() {
		BeforeEach(func() {
			configFile = "non-existing-folder/config.yaml"
		})

		It("fails with an error", func() {
			Expect(err.Error()).To(ContainSubstring("no such file or directory"))
		})
	})

	When("given an invalid config file", func() {
		BeforeEach(func() {
			configFile = invalidConfigPath
		})

		It("fails with an error", func() {
			Expect(err.Error()).To(ContainSubstring("value is not allowed in this context"))
		})
	})
})
