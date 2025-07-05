package git_test

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/AntonKosov/git-backups/internal/cmd"
	"github.com/AntonKosov/git-backups/internal/git"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Git tests", func() {
	const (
		firstCommitArchive  = "./test_data/first_commit.zip"
		secondCommitArchive = "./test_data/second_commit.zip"
		firstCommitID       = "07b22a94cd460978cb91aec6c1b8967973996f3e"
		secondCommitID      = "c95f9eee1bf80bd9ffa198efd239599365cefa40"
	)

	var (
		ctx        context.Context
		err        error
		sourcePath string
		targetPath string
	)

	verifyID := func(expected string) {
		var stdout strings.Builder
		command := exec.Command("git", "-C", targetPath, "rev-list", "HEAD", "-1")
		command.Stdout = &stdout
		err := command.Run()

		Expect(err).NotTo(HaveOccurred())
		Expect(stdout.String()).To(Equal(expected + "\n"))
	}

	clearSource := func() {
		command := exec.Command("rm", "-rf", sourcePath+"/*")
		err := command.Run()
		Expect(err).NotTo(HaveOccurred())
	}

	unzipArchive := func(archive string) {
		clearSource()

		err = cmd.Execute(ctx, "unzip", archive, "-d", sourcePath)
		Expect(err).NotTo(HaveOccurred())
	}

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		sourcePath, err = os.MkdirTemp(os.TempDir(), "source_*")
		Expect(err).NotTo(HaveOccurred())

		targetPath, err = os.MkdirTemp(os.TempDir(), "target_*")
		Expect(err).NotTo(HaveOccurred())

		unzipArchive(firstCommitArchive)
	})

	AfterEach(func() {
		if sourcePath != "" {
			_ = os.RemoveAll(sourcePath)
			sourcePath = ""
		}

		if targetPath != "" {
			_ = os.RemoveAll(targetPath)
			targetPath = ""
		}
	})

	Context("Clone", func() {
		var source string

		BeforeEach(func() {
			source = sourcePath
		})

		JustBeforeEach(func() {
			err = git.Clone(ctx, source, targetPath)
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
	})

	Context("Fetch", func() {
		BeforeEach(func() {
			unzipArchive(secondCommitArchive)
			err := git.Clone(ctx, sourcePath, targetPath)
			Expect(err).NotTo(HaveOccurred())
		})

		JustBeforeEach(func() {
			err = git.Fetch(ctx, targetPath)
		})

		FIt("does not return an error", func() {
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
				Expect(err).To(ContainSubstring(`git failed with "fatal: not a git repository:`))
			})
		})
	})
})
