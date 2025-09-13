package launcher_test

import (
	"context"
	"errors"

	"github.com/AntonKosov/git-backups/internal/config"
	"github.com/AntonKosov/git-backups/internal/github"
	"github.com/AntonKosov/git-backups/internal/launcher"
	"github.com/AntonKosov/git-backups/internal/launcher/launcherfakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Launcher tests", func() {
	var (
		conf = config.Config{
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
		fakeBackupService *launcherfakes.FakeBackupService
		fakeReaderService *launcherfakes.FakeReaderService
		err               error
	)

	BeforeEach(func() {
		fakeBackupService = &launcherfakes.FakeBackupService{}
		fakeReaderService = &launcherfakes.FakeReaderService{}

		setReturnCall := func(callCount int, repos []github.Repo) {
			fakeReaderService.AllReposReturnsOnCall(callCount, func(yield func(github.Repo, error) bool) {
				for _, repo := range repos {
					if !yield(repo, nil) {
						return
					}
				}
			})
		}
		setReturnCall(0, []github.Repo{
			{
				Name:     "repo_name_1",
				Owner:    "GH_Username1",
				CloneURL: "https://github.com/GH_Username1/repo_name_1.git",
			},
			{
				Name:     "repo_name_2",
				Owner:    "GH_Username1",
				CloneURL: "https://github.com/GH_Username1/repo_name_2.git",
			},
			{
				Name:     "repo_name_3",
				Owner:    "GH_Username1",
				CloneURL: "https://github.com/GH_Username1/repo_name_3.git",
			},
		})
		setReturnCall(1, []github.Repo{
			{
				Name:     "repo_name_4",
				Owner:    "GH_Username2",
				CloneURL: "https://github.com/GH_Username2/repo_name_4.git",
			},
			{
				Name:     "repo_name_5",
				Owner:    "GH_Username2",
				CloneURL: "https://github.com/GH_Username2/repo_name_5.git",
			},
			{
				Name:     "repo_name_6",
				Owner:    "GH_Username2",
				CloneURL: "https://github.com/GH_Username2/repo_name_6.git",
			},
			{
				Name:     "non_included_repo_name",
				Owner:    "GH_Username2",
				CloneURL: "https://github.com/GH_Username2/non_included_repo_name.git",
			},
		})
		setReturnCall(2, []github.Repo{
			{
				Name:     "repo_name_7",
				Owner:    "GH_Username3",
				CloneURL: "https://github.com/GH_Username3/repo_name_7.git",
			},
			{
				Name:     "repo_name_8",
				Owner:    "GH_Username3",
				CloneURL: "https://github.com/GH_Username3/repo_name_8.git",
			},
		})
		setReturnCall(3, []github.Repo{
			{
				Name:     "repo_name_9",
				Owner:    "GH_Username4",
				CloneURL: "https://github.com/GH_Username4/repo_name_9.git",
			},
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
		verify := func(idx int, expectedURL, expectedPath string) {
			_, url, path := fakeBackupService.RunArgsForCall(idx)
			Expect(url).To(Equal(expectedURL))
			Expect(path).To(Equal(expectedPath))
		}

		verify(0, "https://github.com/Username1/repo_name_1.git", "/home/user/git_backup/folder_name/repo_folder_name_1")
		verify(1, "https://github.com/Username2/repo_name_2.git", "/home/user/git_backup/folder_name/repo_folder_name_2")
		verify(2, "https://github.com/Username3/repo_name_3.git", "/home/user/git_backup/folder_name_2/repo_folder_name_3")
		verify(3, "https://github.com/Username4/repo_name_4.git", "/home/user/git_backup/folder_name_2/repo_folder_name_4")
		verify(4, "https://oauth2:GH_XXX@github.com/GH_Username1/repo_name_1.git", "/home/user/git_backup/folder_name_3/GH_Username1/repo_name_1")
		verify(5, "https://oauth2:GH_XXX@github.com/GH_Username1/repo_name_2.git", "/home/user/git_backup/folder_name_3/GH_Username1/repo_name_2")
		verify(6, "https://oauth2:GH2_XXX@github.com/GH_Username2/repo_name_4.git", "/home/user/git_backup/folder_name_4/GH_Username2/repo_name_4")
		verify(7, "https://oauth2:GH2_XXX@github.com/GH_Username2/repo_name_5.git", "/home/user/git_backup/folder_name_4/GH_Username2/repo_name_5")
		verify(8, "https://oauth2:GH3_XXX@github.com/GH_Username3/repo_name_7.git", "/home/user/git_backup/folder_name_5/GH_Username3/repo_name_7")
		verify(9, "https://oauth2:GH4_XXX@github.com/GH_Username4/repo_name_9.git", "/home/user/git_backup/folder_name_6/GH_Username4/repo_name_9")
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

	When("include contains a repo name with different capitalization", func() {
		BeforeEach(func() {
			conf.Repositories.GitHub[0].Include[0] = "rEpO_NaMe_1"
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("includes the repo", func() {
			_, url, path := fakeBackupService.RunArgsForCall(4)
			Expect(url).To(Equal("https://oauth2:GH_XXX@github.com/GH_Username1/repo_name_1.git"))
			Expect(path).To(Equal("/home/user/git_backup/folder_name_3/GH_Username1/repo_name_1"))
		})
	})

	When("exclude contains a repo name with different capitalization", func() {
		BeforeEach(func() {
			conf.Repositories.GitHub[0].Exclude[0] = "rEpO_NaMe_3"
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not include the repo", func() {
			_, url, path := fakeBackupService.RunArgsForCall(6)
			Expect(url).To(Equal("https://oauth2:GH2_XXX@github.com/GH_Username2/repo_name_4.git"))
			Expect(path).To(Equal("/home/user/git_backup/folder_name_4/GH_Username2/repo_name_4"))
		})
	})

	When("include is empty", func() {
		BeforeEach(func() {
			conf.Repositories.GitHub[0].Include = make([]string, 0)
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not clone repositories from the first GitHub profile", func() {
			Expect(fakeBackupService.RunCallCount()).To(Equal(8))
			_, url, path := fakeBackupService.RunArgsForCall(4)
			Expect(url).To(Equal("https://oauth2:GH2_XXX@github.com/GH_Username2/repo_name_4.git"))
			Expect(path).To(Equal("/home/user/git_backup/folder_name_4/GH_Username2/repo_name_4"))
		})
	})

	When("generic context is canceled", func() {
		BeforeEach(func() {
			numCalls := 0
			fakeBackupService.RunStub = func(context.Context, string, string) error {
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
			fakeBackupService.RunStub = func(context.Context, string, string) error {
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
