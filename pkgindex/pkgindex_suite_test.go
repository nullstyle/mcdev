package pkgindex_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPkgindex(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pkgindex Suite")
}
