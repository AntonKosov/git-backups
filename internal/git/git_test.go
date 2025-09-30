package git_test

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AntonKosov/git-backups/internal/cmd"
	"github.com/AntonKosov/git-backups/internal/git"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Git tests", func() {
	const (
		firstCommitArchive  = "../../test/data/first_commit.zip"
		secondCommitArchive = "../../test/data/second_commit.zip"
		firstCommitID       = "07b22a94cd460978cb91aec6c1b8967973996f3e"
		secondCommitID      = "c95f9eee1bf80bd9ffa198efd239599365cefa40"
	)

	var (
		err        error
		sourcePath string
		targetPath string
		worker     git.Git
	)

	mkdir := func(name string) {
		Expect(os.Mkdir(name, os.ModeDir|os.ModePerm)).NotTo(HaveOccurred())
	}

	rmdir := func(name string) {
		Expect(os.RemoveAll(name)).NotTo(HaveOccurred())
	}

	mkdirTemp := func(prefix string) string {
		tmpDir, defined := os.LookupEnv("TMP_DIR")
		Expect(defined).To(BeTrue())
		folderName := fmt.Sprintf("%v/%v_%v", tmpDir, prefix, time.Now().UnixNano())
		mkdir(folderName)

		return folderName
	}

	verifyID := func(expected string) {
		var output strings.Builder
		err := cmd.Execute(
			ctx,
			"git",
			cmd.WithArguments("-C", targetPath, "rev-list", "FETCH_HEAD", "-1"),
			cmd.WithStdoutWriter(&output),
		)
		if err != nil {
			if !strings.Contains(err.Error(), "ambiguous argument 'FETCH_HEAD'") {
				Fail("Unexpected error " + err.Error())
			}
			err = cmd.Execute(
				ctx,
				"git",
				cmd.WithArguments("-C", targetPath, "rev-list", "HEAD", "-1"),
				cmd.WithStdoutWriter(&output),
			)
			Expect(err).NotTo(HaveOccurred())
		}

		Expect(output.String()).To(Equal(expected + "\n"))
	}

	clearSource := func() {
		rmdir(sourcePath)
		mkdir(sourcePath)
	}

	unzipArchiveToSource := func(archive string) {
		clearSource()

		err := cmd.Execute(ctx, "unzip", cmd.WithArguments(archive, "-d", sourcePath))
		Expect(err).NotTo(HaveOccurred())
	}

	BeforeEach(func() {
		sourcePath = mkdirTemp("source")
		targetPath = mkdirTemp("target")

		unzipArchiveToSource(firstCommitArchive)
	})

	AfterEach(func() {
		if sourcePath != "" {
			rmdir(sourcePath)
			sourcePath = ""
		}

		if targetPath != "" {
			rmdir(targetPath)
			targetPath = ""
		}
	})

	Context("Clone", func() {
		var (
			source        string
			privateSSHKey *string
		)

		BeforeEach(func() {
			source = sourcePath
			privateSSHKey = nil
		})

		JustBeforeEach(func() {
			err = worker.Clone(ctx, source, targetPath, privateSSHKey)
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("has correct commit ID", func() {
			verifyID(firstCommitID)
		})

		When("source is unavailable", func() {
			BeforeEach(func() {
				source += "/missing_path"
			})

			It("returns an error", func() {
				Expect(err.Error()).To(ContainSubstring("missing_path' does not exist"))
			})
		})

		When("a private SSH key is provided", func() {
			BeforeEach(func() {
				path := "/path/to/private/ssh/key"
				privateSSHKey = &path
			})

			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Context("Fetch", func() {
		var privateSSHKey *string

		BeforeEach(func() {
			privateSSHKey = nil
			err := worker.Clone(ctx, sourcePath, targetPath, privateSSHKey)
			Expect(err).NotTo(HaveOccurred())
			unzipArchiveToSource(secondCommitArchive)
		})

		JustBeforeEach(func() {
			err = worker.Fetch(ctx, targetPath, nil)
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("has correct commit ID", func() {
			verifyID(secondCommitID)
		})

		When("source is unavailable", func() {
			BeforeEach(func() {
				clearSource()
			})

			It("returns an error", func() {
				Expect(err.Error()).To(ContainSubstring(`fatal: Could not read from remote repository.`))
			})
		})

		When("a private SSH key is provided", func() {
			BeforeEach(func() {
				path := "/path/to/private/ssh/key"
				privateSSHKey = &path
			})

			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Context("GetRemoteURL", func() {
		var (
			remoteURL string
			path      string
		)

		BeforeEach(func() {
			err := worker.Clone(ctx, sourcePath, targetPath, nil)
			Expect(err).NotTo(HaveOccurred())
			path = targetPath
		})

		JustBeforeEach(func() {
			remoteURL, err = worker.GetRemoteURL(ctx, path)
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns correct remote URL", func() {
			Expect(remoteURL).To(Equal(sourcePath))
		})

		When("command returns an error", func() {
			BeforeEach(func() {
				path = "./non/existing/path"
			})

			It("returns an error", func() {
				Expect(err.Error()).To(ContainSubstring("No such file or directory"))
			})
		})
	})

	Context("SetRemoteURL", func() {
		const newURL = "https://new_url.com"
		var (
			remoteURL string
			path      string
		)

		BeforeEach(func() {
			err := worker.Clone(ctx, sourcePath, targetPath, nil)
			Expect(err).NotTo(HaveOccurred())
			path = targetPath
		})

		JustBeforeEach(func() {
			err = worker.SetRemoteURL(ctx, path, newURL)
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns correct remote URL", func() {
			remoteURL, err = worker.GetRemoteURL(ctx, targetPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(remoteURL).To(Equal(newURL))
		})

		When("command returns an error", func() {
			BeforeEach(func() {
				path = "./non/existing/path"
			})

			It("returns an error", func() {
				Expect(err.Error()).To(ContainSubstring("No such file or directory"))
			})
		})
	})
})
