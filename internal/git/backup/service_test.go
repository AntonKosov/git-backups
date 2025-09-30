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
		privateSSHKey *string
		fakeGit       *backupfakes.FakeGit
		service       backup.Service
		err           error
	)

	BeforeEach(func() {
		privateSSHKey = nil
		fakeGit = &backupfakes.FakeGit{}
		service = backup.NewService(fakeGit)
	})

	JustBeforeEach(func() {
		err = service.Run(ctx, sourceURL, missingFolder, privateSSHKey)
	})

	It("does not return an error", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	It("clones with correct arguments", func() {
		Expect(fakeGit.CloneCallCount()).To(Equal(1))
		_, url, path, privateSSHKey := fakeGit.CloneArgsForCall(0)
		Expect(url).To(Equal(sourceURL))
		Expect(path).To(Equal(missingFolder))
		Expect(privateSSHKey).To(BeNil())
	})

	When("clone returns an error", func() {
		BeforeEach(func() {
			fakeGit.CloneReturns(errors.New("something went wrong"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("something went wrong"))
		})
	})

	When("a private SSH key is provided", func() {
		BeforeEach(func() {
			key := "/path/to/ssh/key"
			privateSSHKey = &key
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("clones with correct arguments", func() {
			Expect(fakeGit.CloneCallCount()).To(Equal(1))
			_, url, path, privateSSHKey := fakeGit.CloneArgsForCall(0)
			Expect(url).To(Equal(sourceURL))
			Expect(path).To(Equal(missingFolder))
			Expect(*privateSSHKey).To(Equal("/path/to/ssh/key"))
		})
	})

	When("target folder exists", func() {
		JustBeforeEach(func() {
			err = service.Run(ctx, sourceURL, targetFolder, privateSSHKey)
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("fetches with correct arguments", func() {
			Expect(fakeGit.FetchCallCount()).To(Equal(1))
			_, path, privateSSHKey := fakeGit.FetchArgsForCall(0)
			Expect(path).To(Equal(targetFolder))
			Expect(privateSSHKey).To(BeNil())
		})

		It("doesn't change the remote URL", func() {
			Expect(fakeGit.SetRemoteURLCallCount()).To(Equal(1))
		})

		When("a private SSH key is provided", func() {
			BeforeEach(func() {
				key := "/path/to/ssh/key"
				privateSSHKey = &key
			})

			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("fetches with correct arguments", func() {
				Expect(fakeGit.FetchCallCount()).To(Equal(1))
				_, path, privateSSHKey := fakeGit.FetchArgsForCall(0)
				Expect(path).To(Equal(targetFolder))
				Expect(*privateSSHKey).To(Equal("/path/to/ssh/key"))
			})
		})

		When("fetch returns an error", func() {
			BeforeEach(func() {
				fakeGit.FetchReturns(errors.New("something went wrong"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("something went wrong"))
			})
		})

		When("get remote URL returns an error", func() {
			BeforeEach(func() {
				fakeGit.GetRemoteURLReturns("", errors.New("something went wrong"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("something went wrong"))
			})
		})

		When("set remote URL returns an error", func() {
			BeforeEach(func() {
				fakeGit.SetRemoteURLReturns(errors.New("something went wrong"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("something went wrong"))
			})
		})

		When("URLs do not match", func() {
			BeforeEach(func() {
				fakeGit.GetRemoteURLReturns("https://www.different_url.com", nil)
			})

			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("calls SetRemoteURL with correct parameters", func() {
				Expect(fakeGit.SetRemoteURLCallCount()).To(Equal(1))
				_, path, url := fakeGit.SetRemoteURLArgsForCall(0)
				Expect(path).To(Equal(targetFolder))
				Expect(url).To(Equal(sourceURL))
			})
		})
	})
})
