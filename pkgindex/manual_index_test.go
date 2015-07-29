package pkgindex

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ManualIndex", func() {
	var subject *ManualIndex

	BeforeEach(func() {
		subject = &ManualIndex{}
	})

	Describe("Add", func() {
		It("adds all the provided packages to the indexes internal array", func() {
			subject.Add("a", "b")
			Expect(subject.packages).To(HaveKey("a"))
			Expect(subject.packages).To(HaveKey("b"))
		})

		It("doesn't add duplicates", func() {
			subject.Add("a")
			subject.Add("a")
			Expect(subject.packages).To(HaveLen(1))
			subject.packages = nil
			subject.Add("a", "a")
			Expect(subject.packages).To(HaveLen(1))
		})
	})

	Describe("Index", func() {
		It("creates the suffix array", func() {
			Expect(subject.index).To(BeNil())
			subject.Index()
			Expect(subject.index).NotTo(BeNil())
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
			subject.Add(packages...)
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

		// TODO: put it own file, use regular benchmarking suite
		//
		// indexWithRandomPackages := func(n int) *ManualIndex {
		// 	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		// 	packages := make([]string, n)
		//
		// 	for i := 0; i < n; i++ {
		// 		pkgLength := rand.Intn(9) + 1
		// 		pkg := make([]rune, pkgLength)
		// 		for i := range pkg {
		// 			pkg[i] = letters[rand.Intn(len(letters))]
		// 		}
		// 		packages[i] = "github.org/" + string(pkg)
		// 	}
		//
		// 	result := &ManualIndex{}
		// 	err := result.Add(packages...)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	return result
		// }
		//
		// measureSearches := func(subject *ManualIndex) {
		// 	Measure("empty search", func(b Benchmarker) {
		// 		b.Time("runtime", func() {
		// 			_, err := subject.Search("")
		// 			Expect(err).To(BeNil())
		// 		})
		// 	}, 1000)
		//
		// 	Measure("query with many matches", func(b Benchmarker) {
		// 		b.Time("runtime", func() {
		// 			_, err := subject.Search("github.org")
		// 			Expect(err).To(BeNil())
		// 		})
		// 	}, 1000)
		//
		// 	Measure("query with fewer matches", func(b Benchmarker) {
		// 		b.Time("runtime", func() {
		// 			_, err := subject.Search("github.org/aa")
		// 			Expect(err).To(BeNil())
		// 		})
		// 	}, 1000)
		// }
		//
		// Context("with 100 entries", func() {
		// 	measureSearches(indexWithRandomPackages(100))
		// })
		//
		// Context("with 1000 entries", func() {
		// 	measureSearches(indexWithRandomPackages(1000))
		// })
	})

})
