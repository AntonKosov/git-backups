package slice_test

import (
	"strconv"

	"github.com/AntonKosov/git-backups/internal/slice"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Lookup", func(input []int, expectedOutput map[byte]string) {
	output := slice.Lookup(input, func(v int) (byte, string) { return byte(v), strconv.Itoa(v) })
	Expect(output).To(Equal(expectedOutput))
},
	Entry("nil input", nil, map[byte]string{}),
	Entry("empty slice", []int{}, map[byte]string{}),
	Entry("one item", []int{5}, map[byte]string{5: "5"}),
	Entry("three items", []int{9, 1, 7}, map[byte]string{9: "9", 1: "1", 7: "7"}),
)
