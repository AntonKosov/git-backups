package git_test

import (
	"errors"

	"github.com/AntonKosov/git-backups/internal/git"
	"github.com/AntonKosov/git-backups/internal/git/gitfakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fetcher tests", func() {
	const (
		sourceURL     = "https://www.abc.com"
		targetFolder  = "./test_data"
		missingFolder = "missing_folder"
	)

	var (
		fakeGitWorker *gitfakes.FakeGitWorker
		fetcher       git.Fetcher
		err           error
	)

	BeforeEach(func() {
		fakeGitWorker = &gitfakes.FakeGitWorker{}
		fetcher = git.NewFetcher(fakeGitWorker)
	})

	JustBeforeEach(func() {
		err = fetcher.Run(ctx, sourceURL, missingFolder)
	})

	It("does not return an error", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	It("clones with correct arguments", func() {
		Expect(fakeGitWorker.CloneCallCount()).To(Equal(1))
		_, url, path := fakeGitWorker.CloneArgsForCall(0)
		Expect(url).To(Equal(sourceURL))
		Expect(path).To(Equal(missingFolder))
	})

	When("clone returns an error", func() {
		BeforeEach(func() {
			fakeGitWorker.CloneReturns(errors.New("something went wrong"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("something went wrong"))
		})
	})

	When("target folder exists", func() {
		JustBeforeEach(func() {
			err = fetcher.Run(ctx, sourceURL, targetFolder)
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("fetches with correct arguments", func() {
			Expect(fakeGitWorker.FetchCallCount()).To(Equal(1))
			_, path := fakeGitWorker.FetchArgsForCall(0)
			Expect(path).To(Equal(targetFolder))
		})

		When("fetch returns an error", func() {
			BeforeEach(func() {
				fakeGitWorker.FetchReturns(errors.New("something went wrong"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("something went wrong"))
			})
		})
	})
})
