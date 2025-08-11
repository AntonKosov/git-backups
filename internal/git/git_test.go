package git_test

import (
	"fmt"
	"os"
	"os/exec"
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
		var stdout strings.Builder
		command := exec.Command("git", "-C", targetPath, "rev-list", "HEAD", "-1")
		command.Stdout = &stdout
		Expect(command.Run()).NotTo(HaveOccurred())
		Expect(stdout.String()).To(Equal(expected + "\n"))
	}

	clearSource := func() {
		rmdir(sourcePath)
		mkdir(sourcePath)
	}

	unzipArchive := func(archive string) {
		clearSource()

		err := cmd.Execute(ctx, "unzip", archive, "-d", sourcePath)
		Expect(err).NotTo(HaveOccurred())
	}

	BeforeEach(func() {
		sourcePath = mkdirTemp("source")
		targetPath = mkdirTemp("target")

		unzipArchive(firstCommitArchive)
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
		var source string

		BeforeEach(func() {
			source = sourcePath
		})

		JustBeforeEach(func() {
			err = worker.Clone(ctx, source, targetPath)
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
			err := worker.Clone(ctx, sourcePath, targetPath)
			Expect(err).NotTo(HaveOccurred())
		})

		JustBeforeEach(func() {
			err = worker.Fetch(ctx, targetPath)
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
	})
})
