package backup_test

import (
	"errors"

	"github.com/AntonKosov/git-backups/internal/git/backup"
	"github.com/AntonKosov/git-backups/internal/git/backup/backupfakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service tests", func() {
	const (
		sourceURL     = "https://www.abc.com"
		targetFolder  = "../../../test/data"
		missingFolder = "missing_folder"
	)

	var (
		fakeGit *backupfakes.FakeGit
		service backup.Service
		err     error
	)

	BeforeEach(func() {
		fakeGit = &backupfakes.FakeGit{}
		service = backup.NewService(fakeGit)
	})

	JustBeforeEach(func() {
		err = service.Run(ctx, sourceURL, missingFolder)
	})

	It("does not return an error", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	It("clones with correct arguments", func() {
		Expect(fakeGit.CloneCallCount()).To(Equal(1))
		_, url, path := fakeGit.CloneArgsForCall(0)
		Expect(url).To(Equal(sourceURL))
		Expect(path).To(Equal(missingFolder))
	})

	When("clone returns an error", func() {
		BeforeEach(func() {
			fakeGit.CloneReturns(errors.New("something went wrong"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("something went wrong"))
		})
	})

	When("target folder exists", func() {
		JustBeforeEach(func() {
			err = service.Run(ctx, sourceURL, targetFolder)
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("fetches with correct arguments", func() {
			Expect(fakeGit.FetchCallCount()).To(Equal(1))
			_, path := fakeGit.FetchArgsForCall(0)
			Expect(path).To(Equal(targetFolder))
		})

		When("fetch returns an error", func() {
			BeforeEach(func() {
				fakeGit.FetchReturns(errors.New("something went wrong"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("something went wrong"))
			})
		})
	})
})
