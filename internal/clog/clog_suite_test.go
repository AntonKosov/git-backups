package clog_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestClog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Clog Suite")
}
