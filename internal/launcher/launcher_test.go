package launcher_test

import (
	//"github.com/AntonKosov/git-backups/internal/launcher"
	"errors"

	"github.com/AntonKosov/git-backups/internal/config"
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
				/*GitHub: []config.GitHubRepo{
					{
						Name:       "profile name 3",
						RootFolder: "/home/user/git_backup/folder_name_3",
						Token:      "GH_XXX",
						Include: []string{
							"repo_name_1",
							"repo_name_2",
						},
						Exclude: []string{
							"repo_name_3",
						},
					},
					{
						Name:       "profile name 4",
						RootFolder: "/home/user/git_backup/folder_name_4",
						Token:      "GH2_XXX",
						Include: []string{
							"repo_name_4",
							"repo_name_5",
						},
						Exclude: []string{
							"repo_name_6",
						},
					},
				},*/
			},
		}
		fakeBackupService *launcherfakes.FakeBackupService
		err               error
	)

	BeforeEach(func() {
		fakeBackupService = &launcherfakes.FakeBackupService{}
	})

	JustBeforeEach(func() {
		err = launcher.Run(ctx, conf, fakeBackupService)
	})

	It("does not return an error", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	It("is called with correct arguments", func() {
		Expect(fakeBackupService.RunCallCount()).To(Equal(4))
		verify := func(idx int, expectedURL, expectedPath string) {
			_, url, path := fakeBackupService.RunArgsForCall(idx)
			Expect(url).To(Equal(expectedURL))
			Expect(path).To(Equal(expectedPath))
		}

		verify(0, "https://github.com/Username1/repo_name_1.git", "/home/user/git_backup/folder_name/repo_folder_name_1")
		verify(1, "https://github.com/Username2/repo_name_2.git", "/home/user/git_backup/folder_name/repo_folder_name_2")
		verify(2, "https://github.com/Username3/repo_name_3.git", "/home/user/git_backup/folder_name_2/repo_folder_name_3")
		verify(3, "https://github.com/Username4/repo_name_4.git", "/home/user/git_backup/folder_name_2/repo_folder_name_4")
	})

	When("backup service returns an error", func() {
		BeforeEach(func() {
			fakeBackupService.RunReturns(errors.New("something went wrong"))
		})

		It("returns an error", func() {
			Expect(err).To(MatchError("something went wrong"))
		})
	})
})
