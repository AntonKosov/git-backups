package internal_test

import (
	"github.com/AntonKosov/git-backups/internal"
	"github.com/AntonKosov/git-backups/internal/internalfakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sum tests", func() {
	var (
		fakeTransformer internalfakes.FakeTransformer
		sum             int
	)

	BeforeEach(func() {
		fakeTransformer = internalfakes.FakeTransformer{}
		fakeTransformer.TransformStub = func(v int) int { return v * 10 }
	})

	JustBeforeEach(func() {
		sum = internal.Sum(2, 3, &fakeTransformer)
	})

	It("returns the correct sum", func() {
		Expect(sum).To(Equal(50))
	})
})
