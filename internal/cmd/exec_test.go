package cmd_test

import (
	"context"
	"errors"
	"strings"

	"github.com/AntonKosov/git-backups/internal/cmd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exec tests", func() {
	var (
		err            error
		executableApp  string
		commandOptions []cmd.Option
	)

	BeforeEach(func() {
		executableApp = "ls"
		commandOptions = nil
	})

	JustBeforeEach(func() {
		err = cmd.Execute(context.Background(), executableApp, commandOptions...)
	})

	It("doesn't return an error", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	When("app has arguments", func() {
		BeforeEach(func() {
			commandOptions = append(commandOptions, cmd.WithArguments("-a"))
		})

		It("doesn't return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("app has incorrect arguments", func() {
		BeforeEach(func() {
			commandOptions = append(commandOptions, cmd.WithArguments("non-existing-folder"))
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

	When("app returns correct output", func() {
		var stdout strings.Builder

		BeforeEach(func() {
			stdout = strings.Builder{}
			commandOptions = append(commandOptions, cmd.WithStdoutWriter(&stdout))
		})

		It("doesn't return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns correct output", func() {
			Expect(stdout.String()).To(ContainSubstring("exec_test.go"))
		})
	})

	When("there are env variables", func() {
		var stdout strings.Builder

		BeforeEach(func() {
			stdout = strings.Builder{}
			executableApp = "env"
			commandOptions = append(
				commandOptions,
				cmd.WithStdoutWriter(&stdout),
				cmd.WithEnvVariables(`var1=1234`, `var2="value with spaces"`),
			)
		})

		It("doesn't return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns correct output", func() {
			Expect(stdout.String()).To(Equal("var1=1234\nvar2=\"value with spaces\"\n"))
		})
	})

	When("an option is nil", func() {
		BeforeEach(func() {
			commandOptions = append(commandOptions, nil)
		})

		It("doesn't return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
