package cmd_test

import (
	"context"
	"errors"

	"github.com/AntonKosov/git-backups/internal/cmd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exec tests", func() {
	var (
		err  error
		args []string
	)

	BeforeEach(func() {
		args = nil
	})

	JustBeforeEach(func() {
		err = cmd.Execute(context.Background(), "ls", args...)
	})

	It("doesn't return an error", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	When("app has arguments", func() {
		BeforeEach(func() {
			args = []string{"-a"}
		})

		It("doesn't return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("app has incorrect arguments", func() {
		BeforeEach(func() {
			args = []string{"non-existing-folder"}
		})

		It("returns an error", func() {
			Expect(err).To(HaveOccurred())
			var cmdErr cmd.CommandError
			Expect(errors.As(err, &cmdErr)).To(BeTrue())
			Expect(cmdErr).To(Equal(cmd.CommandError{
				Name: "ls",
				Args: []string{
					"non-existing-folder",
				},
				Err: "ls: cannot access 'non-existing-folder': No such file or directory\n",
			}))
		})
	})
})
