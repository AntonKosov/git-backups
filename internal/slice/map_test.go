package slice_test

import (
	"strconv"

	"github.com/AntonKosov/git-backups/internal/slice"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Map", func(input []int, expectedOutput []string) {
	output := slice.Map(input, func(v int) string { return strconv.Itoa(v) })
	Expect(output).To(Equal(expectedOutput))
},
	Entry("nil input", nil, nil),
	Entry("empty slice", []int{}, []string{}),
	Entry("one item", []int{5}, []string{"5"}),
	Entry("three items", []int{9, 1, 7}, []string{"9", "1", "7"}),
)
