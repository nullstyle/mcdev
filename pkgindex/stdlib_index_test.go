package pkgindex

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StdLibIndex", func() {
	var subject *StdLibIndex

	BeforeEach(func() {
		subject = &StdLibIndex{}
		err := subject.Index()
		if err != nil {
			Fail("failed to index std library")
		}
	})

	AfterEach(func() {
		//TODO: clear index
	})

	Describe("Search", func() {
		It("finds all 'encoding/base' packages", func() {
			results, _ := subject.Search("encoding/base")
			Expect(results).To(HaveLen(2))
			Expect(results).To(ContainElement("encoding/base32"))
			Expect(results).To(ContainElement("encoding/base64"))
		})

		It("finds 'text/template/parse' packages", func() {
			results, _ := subject.Search("tempp")
			Expect(results).To(HaveLen(1))
			Expect(results).To(ContainElement("text/template/parse"))
		})

		// TODO: add some benchmarks
	})
})
