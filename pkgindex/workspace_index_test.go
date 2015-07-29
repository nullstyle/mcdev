package pkgindex

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WorkspaceIndex", func() {
	var subject *WorkspaceIndex
	var dir string

	BeforeSuite(func() {
		var err error
		dir, err = ioutil.TempDir("", "mcdev-workspace-index")
		if err != nil {
			Fail("could not create tmpdir")
		}

		touch := func(path string) {
			path = filepath.Join(dir, "src", path)
			err := os.MkdirAll(path, 0755)
			if err != nil {
				Fail("Couldn't create package dir")
			}

			err = ioutil.WriteFile(
				filepath.Join(path, "main.go"),
				[]byte("package main"),
				0755,
			)
			if err != nil {
				Fail("Couldn't create dummy .go file")
			}
		}

		// dirs to be indexed
		touch("github.org/nullstyle/mcdev")
		touch("golang.org/x/net/context")
		// dirs to be ignored
		touch("github.org/.hidden")

		subject = &WorkspaceIndex{Dir: dir}
	})

	AfterSuite(func() {
		os.RemoveAll(dir)
	})

	Describe("Index", func() {
		BeforeEach(func() {
			err := subject.Index()
			Expect(err).To(BeNil())
		})

		It("adds all directories underneath the workspace into the index", func() {
			Expect(subject.idx.packages).To(HaveLen(2))
			Expect(subject.idx.packages).NotTo(HaveKey("github.org/nullstyle/mcdev"))
			Expect(subject.idx.packages).NotTo(HaveKey("golang.org/x/net/context"))
		})

		It("ignores files", func() {
			Expect(subject.idx.packages).NotTo(HaveKey("github.org/nullstyle/mcdev/main.go"))
		})

		It("ignores hidden directories", func() {
			Expect(subject.idx.packages).NotTo(HaveKey("github.org/.hidden"))
		})
	})
})
