package pkgindex

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CompoundIndex", func() {
	var subject *CompoundIndex

	BeforeEach(func() {
		subject = &CompoundIndex{}
	})

	Describe("Add", func() {
		It("add the child index into the children slice", func() {
			cidx := &ManualIndex{}
			Expect(subject.children).NotTo(ContainElement(cidx))
			subject.Add(cidx)
			Expect(subject.children).To(ContainElement(cidx))
		})
	})

	Describe("Index", func() {
		It("calls Index on each child", func() {
			one := &ManualIndex{}
			two := &ManualIndex{}

			subject.Add(one)
			subject.Add(two)

			Expect(one.index).To(BeNil())
			Expect(two.index).To(BeNil())

			subject.Index()

			Expect(one.index).NotTo(BeNil())
			Expect(two.index).NotTo(BeNil())
		})
	})

	Describe("Search", func() {
		packages := []string{
			"github.org/nullstyle/mcdev",
			"github.org/nullstyle/mcdev/imports",
			"github.org/nullstyle/mcdev/pkgwatch",
			"golang.org/x/net/context",
		}

		BeforeEach(func() {
			one := &ManualIndex{}
			two := &ManualIndex{}

			subject.Add(one)
			subject.Add(two)

			one.Add(packages[:2]...)
			two.Add(packages[2:]...)
		})

		Context("the search string is empty", func() {
			It("returns all packages", func() {
				results, err := subject.Search("")
				Expect(err).To(BeNil())
				Expect(results).To(HaveLen(len(packages)))
				// ensure all the records are returned
				for _, pkg := range packages {
					Expect(results).To(ContainElement(pkg))
				}
			})
		})

		Context("the search string 'v'", func() {
			It("returns packages that contain the letter in any position", func() {
				results, err := subject.Search("v")
				Expect(err).To(BeNil())
				Expect(results).To(HaveLen(3))
				Expect(results).To(ContainElement("github.org/nullstyle/mcdev"))
				Expect(results).To(ContainElement("github.org/nullstyle/mcdev/imports"))
				Expect(results).To(ContainElement("github.org/nullstyle/mcdev/pkgwatch"))
				Expect(results).NotTo(ContainElement("golang.org/x/net/context"))
			})
		})

		Context("the search string 'mci'", func() {
			It("returns only the mcdev/imports package", func() {
				results, err := subject.Search("mci")
				Expect(err).To(BeNil())
				Expect(results).To(HaveLen(1))
				Expect(results).To(ContainElement("github.org/nullstyle/mcdev/imports"))
			})
		})

		Context("the search string 'golang'", func() {
			It("returns only the golang.org package", func() {
				results, err := subject.Search("golang")
				Expect(err).To(BeNil())
				Expect(results).To(HaveLen(1))
				Expect(results).To(ContainElement("golang.org/x/net/context"))
			})
		})
	})

})
