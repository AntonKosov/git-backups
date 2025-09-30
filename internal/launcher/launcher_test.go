package launcher_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/AntonKosov/git-backups/internal/config"
	"github.com/AntonKosov/git-backups/internal/github"
	"github.com/AntonKosov/git-backups/internal/launcher"
	"github.com/AntonKosov/git-backups/internal/launcher/launcherfakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Launcher tests", func() {
	var (
		conf              config.Config
		fakeBackupService *launcherfakes.FakeBackupService
		fakeReaderService *launcherfakes.FakeReaderService
		err               error
	)

	var verifyCall = func(idx int, expectedURL, expectedPath string, expectedSSHKey *string) {
		_, url, path, privateSSHKey := fakeBackupService.RunArgsForCall(idx)
		Expect(url).To(Equal(expectedURL))
		Expect(path).To(Equal(expectedPath))
		Expect(privateSSHKey).To(Equal(expectedSSHKey))
	}

	BeforeEach(func() {
		fakeBackupService = &launcherfakes.FakeBackupService{}
		fakeReaderService = &launcherfakes.FakeReaderService{}

		conf = config.Config{
			Profiles: config.Profiles{
				GenericProfiles: []config.GenericProfile{
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
				GitHubProfiles: []config.GitHubProfile{
					{
						Name:        "profile name 3",
						RootFolder:  "/home/user/git_backup/folder_name_3",
						Affiliation: "owner",
						Token:       "GH_XXX",
						Include: []string{
							"repo_name_1",
							"repo_name_2",
							"repo_name_3",
						},
						Exclude: []string{
							"repo_name_3",
						},
					},
					{
						Name:        "profile name 4",
						RootFolder:  "/home/user/git_backup/folder_name_4",
						Affiliation: "owner,collaborator",
						Token:       "GH2_XXX",
						Include: []string{
							"repo_name_4",
							"repo_name_5",
						},
						Exclude: []string{
							"repo_name_6",
						},
					},
					{
						Name:        "profile name 5",
						RootFolder:  "/home/user/git_backup/folder_name_5",
						Affiliation: "owner",
						Token:       "GH3_XXX",
						Exclude: []string{
							"repo_name_8",
						},
					},
					{
						Name:        "profile name 6",
						RootFolder:  "/home/user/git_backup/folder_name_6",
						Affiliation: "owner",
						Token:       "GH4_XXX",
						Include: []string{
							"repo_name_9",
						},
					},
				},
			},
		}

		setReturnCall := func(callCount int, repos []github.Repo) {
			fakeReaderService.AllReposReturnsOnCall(callCount, func(yield func(github.Repo, error) bool) {
				for _, repo := range repos {
					if !yield(repo, nil) {
						return
					}
				}
			})
		}
		ghRepo := func(owner, repoName string) github.Repo {
			return github.Repo{
				Name:     repoName,
				Owner:    owner,
				CloneURL: fmt.Sprintf("https://github.com/%v/%v.git", owner, repoName),
				GitURL:   fmt.Sprintf("git:github.com/%v/%v.git", owner, repoName),
			}
		}
		setReturnCall(0, []github.Repo{
			ghRepo("GH_Username1", "repo_name_1"),
			ghRepo("GH_Username1", "repo_name_2"),
			ghRepo("GH_Username1", "repo_name_3"),
			ghRepo("GH_Username1", "repo_name_3"),
		})
		setReturnCall(1, []github.Repo{
			ghRepo("GH_Username2", "repo_name_4"),
			ghRepo("GH_Username2", "repo_name_5"),
			ghRepo("GH_Username2", "repo_name_6"),
			ghRepo("GH_Username2", "non_included_repo_name"),
		})
		setReturnCall(2, []github.Repo{
			ghRepo("GH_Username3", "repo_name_7"),
			ghRepo("GH_Username3", "repo_name_8"),
		})
		setReturnCall(3, []github.Repo{
			ghRepo("GH_Username4", "repo_name_9"),
		})
	})

	JustBeforeEach(func() {
		err = launcher.Run(ctx, conf, fakeBackupService, fakeReaderService)
	})

	It("does not return an error", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	It("is called with correct arguments", func() {
		Expect(fakeBackupService.RunCallCount()).To(Equal(10))

		verifyCall(0, "https://github.com/Username1/repo_name_1.git", "/home/user/git_backup/folder_name/repo_folder_name_1", nil)
		verifyCall(1, "https://github.com/Username2/repo_name_2.git", "/home/user/git_backup/folder_name/repo_folder_name_2", nil)
		verifyCall(2, "https://github.com/Username3/repo_name_3.git", "/home/user/git_backup/folder_name_2/repo_folder_name_3", nil)
		verifyCall(3, "https://github.com/Username4/repo_name_4.git", "/home/user/git_backup/folder_name_2/repo_folder_name_4", nil)
		verifyCall(4, "https://oauth2:GH_XXX@github.com/GH_Username1/repo_name_1.git", "/home/user/git_backup/folder_name_3/GH_Username1/repo_name_1", nil)
		verifyCall(5, "https://oauth2:GH_XXX@github.com/GH_Username1/repo_name_2.git", "/home/user/git_backup/folder_name_3/GH_Username1/repo_name_2", nil)
		verifyCall(6, "https://oauth2:GH2_XXX@github.com/GH_Username2/repo_name_4.git", "/home/user/git_backup/folder_name_4/GH_Username2/repo_name_4", nil)
		verifyCall(7, "https://oauth2:GH2_XXX@github.com/GH_Username2/repo_name_5.git", "/home/user/git_backup/folder_name_4/GH_Username2/repo_name_5", nil)
		verifyCall(8, "https://oauth2:GH3_XXX@github.com/GH_Username3/repo_name_7.git", "/home/user/git_backup/folder_name_5/GH_Username3/repo_name_7", nil)
		verifyCall(9, "https://oauth2:GH4_XXX@github.com/GH_Username4/repo_name_9.git", "/home/user/git_backup/folder_name_6/GH_Username4/repo_name_9", nil)
	})

	When("generic backup service returns an error", func() {
		BeforeEach(func() {
			fakeBackupService.RunReturns(errors.New("something went wrong"))
		})

		It("returns an error", func() {
			Expect(err.Error()).To(ContainSubstring("failed to backup repository https://github.com/Username1/repo_name_1.git from profile profile name: something went wrong"))
			Expect(err.Error()).To(ContainSubstring("failed to backup repository https://github.com/Username2/repo_name_2.git from profile profile name: something went wrong"))
			Expect(err.Error()).To(ContainSubstring("failed to backup repository https://github.com/Username3/repo_name_3.git from profile profile name 2: something went wrong"))
			Expect(err.Error()).To(ContainSubstring("failed to backup repository https://github.com/Username4/repo_name_4.git from profile profile name 2: something went wrong"))
			Expect(err.Error()).To(ContainSubstring("failed to backup repository https://github.com/GH_Username1/repo_name_1.git from profile profile name 3: something went wrong"))
			Expect(err.Error()).To(ContainSubstring("failed to backup repository https://github.com/GH_Username1/repo_name_2.git from profile profile name 3: something went wrong"))
			Expect(err.Error()).To(ContainSubstring("failed to backup repository https://github.com/GH_Username2/repo_name_4.git from profile profile name 4: something went wrong"))
			Expect(err.Error()).To(ContainSubstring("failed to backup repository https://github.com/GH_Username2/repo_name_5.git from profile profile name 4: something went wrong"))
			Expect(err.Error()).To(ContainSubstring("failed to backup repository https://github.com/GH_Username3/repo_name_7.git from profile profile name 5: something went wrong"))
			Expect(err.Error()).To(ContainSubstring("failed to backup repository https://github.com/GH_Username4/repo_name_9.git from profile profile name 6: something went wrong"))
		})
	})

	When("reader iterator returns an error", func() {
		BeforeEach(func() {
			fakeReaderService.AllReposReturnsOnCall(0, func(yield func(github.Repo, error) bool) {
				yield(github.Repo{}, errors.New("something went wrong"))
			})
		})

		It("returns an error", func() {
			Expect(err.Error()).To(ContainSubstring("failed to read repositories: something went wrong"))
		})
	})

	When("incorrect clone URL is given", func() {
		BeforeEach(func() {
			fakeReaderService.AllReposReturnsOnCall(0, func(yield func(github.Repo, error) bool) {
				yield(github.Repo{
					Name:     "repo_name_1",
					Owner:    "GH_Username1",
					CloneURL: "not_https://github.com/GH_Username1/repo_name_1.git",
				}, nil)
			})
		})

		It("returns an error", func() {
			Expect(err.Error()).To(ContainSubstring("failed to add token to clone URL not_https://github.com/GH_Username1/repo_name_1.git from profile profile name 3: unexpected URL prefix"))
		})
	})

	When("GitHub backup service returns an error", func() {
		BeforeEach(func() {
			fakeBackupService.RunReturnsOnCall(4, errors.New("something went wrong"))
		})

		It("returns an error", func() {
			Expect(err.Error()).To(ContainSubstring("failed to backup repository https://github.com/GH_Username1/repo_name_1.git from profile profile name 3: something went wrong"))
		})
	})

	When("an SSH key is provided", func() {
		var (
			key1 = "/generic/profiles/1/private_key"
			key2 = "/github/profiles/0/private_key"
		)

		BeforeEach(func() {
			genericProfile := &conf.Profiles.GenericProfiles[1]
			genericProfile.PrivateSSHKey = &key1
			genericProfile.Targets[0].URL = "git:github.com/Username3/repo_name_3.git"
			genericProfile.Targets[1].URL = "git:github.com/Username4/repo_name_4.git"

			conf.Profiles.GitHubProfiles[0].PrivateSSHKey = &key2
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("uses the SSH keys", func() {
			Expect(fakeBackupService.RunCallCount()).To(Equal(10))

			verifyCall(0, "https://github.com/Username1/repo_name_1.git", "/home/user/git_backup/folder_name/repo_folder_name_1", nil)
			verifyCall(1, "https://github.com/Username2/repo_name_2.git", "/home/user/git_backup/folder_name/repo_folder_name_2", nil)
			verifyCall(2, "git:github.com/Username3/repo_name_3.git", "/home/user/git_backup/folder_name_2/repo_folder_name_3", &key1)
			verifyCall(3, "git:github.com/Username4/repo_name_4.git", "/home/user/git_backup/folder_name_2/repo_folder_name_4", &key1)
			verifyCall(4, "git:github.com/GH_Username1/repo_name_1.git", "/home/user/git_backup/folder_name_3/GH_Username1/repo_name_1", &key2)
			verifyCall(5, "git:github.com/GH_Username1/repo_name_2.git", "/home/user/git_backup/folder_name_3/GH_Username1/repo_name_2", &key2)
			verifyCall(6, "https://oauth2:GH2_XXX@github.com/GH_Username2/repo_name_4.git", "/home/user/git_backup/folder_name_4/GH_Username2/repo_name_4", nil)
			verifyCall(7, "https://oauth2:GH2_XXX@github.com/GH_Username2/repo_name_5.git", "/home/user/git_backup/folder_name_4/GH_Username2/repo_name_5", nil)
			verifyCall(8, "https://oauth2:GH3_XXX@github.com/GH_Username3/repo_name_7.git", "/home/user/git_backup/folder_name_5/GH_Username3/repo_name_7", nil)
			verifyCall(9, "https://oauth2:GH4_XXX@github.com/GH_Username4/repo_name_9.git", "/home/user/git_backup/folder_name_6/GH_Username4/repo_name_9", nil)
		})
	})

	When("include contains a repo name with different capitalization", func() {
		BeforeEach(func() {
			conf.Profiles.GitHubProfiles[0].Include[0] = "rEpO_NaMe_1"
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("includes the repo", func() {
			_, url, path, privateSSHKey := fakeBackupService.RunArgsForCall(4)
			Expect(url).To(Equal("https://oauth2:GH_XXX@github.com/GH_Username1/repo_name_1.git"))
			Expect(path).To(Equal("/home/user/git_backup/folder_name_3/GH_Username1/repo_name_1"))
			Expect(privateSSHKey).To(BeNil())
		})
	})

	When("exclude contains a repo name with different capitalization", func() {
		BeforeEach(func() {
			conf.Profiles.GitHubProfiles[0].Exclude[0] = "rEpO_NaMe_3"
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not include the repo", func() {
			_, url, path, privateSSHKey := fakeBackupService.RunArgsForCall(6)
			Expect(url).To(Equal("https://oauth2:GH2_XXX@github.com/GH_Username2/repo_name_4.git"))
			Expect(path).To(Equal("/home/user/git_backup/folder_name_4/GH_Username2/repo_name_4"))
			Expect(privateSSHKey).To(BeNil())
		})
	})

	When("include is empty", func() {
		BeforeEach(func() {
			conf.Profiles.GitHubProfiles[0].Include = make([]string, 0)
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not clone repositories from the first GitHub profile", func() {
			Expect(fakeBackupService.RunCallCount()).To(Equal(8))
			_, url, path, privateSSHKey := fakeBackupService.RunArgsForCall(4)
			Expect(url).To(Equal("https://oauth2:GH2_XXX@github.com/GH_Username2/repo_name_4.git"))
			Expect(path).To(Equal("/home/user/git_backup/folder_name_4/GH_Username2/repo_name_4"))
			Expect(privateSSHKey).To(BeNil())
		})
	})

	When("generic context is canceled", func() {
		BeforeEach(func() {
			numCalls := 0
			fakeBackupService.RunStub = func(context.Context, string, string, *string) error {
				numCalls++
				if numCalls == 3 {
					ctxCancel()
				}

				return nil
			}
		})

		It("returns an error", func() {
			Expect(err).To(MatchError(context.Canceled))
		})

		It("was terminated", func() {
			Expect(fakeBackupService.RunCallCount()).To(Equal(3))
		})
	})

	When("GitHub context is canceled", func() {
		BeforeEach(func() {
			numCalls := 0
			fakeBackupService.RunStub = func(context.Context, string, string, *string) error {
				numCalls++
				if numCalls == 6 {
					ctxCancel()
				}

				return nil
			}
		})

		It("returns an error", func() {
			Expect(err).To(MatchError(context.Canceled))
		})

		It("was terminated", func() {
			Expect(fakeBackupService.RunCallCount()).To(Equal(6))
		})
	})
})
