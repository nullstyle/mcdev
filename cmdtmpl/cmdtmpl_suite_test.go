package cmdtmpl_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCmdtmpl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmdtmpl Suite")
}
