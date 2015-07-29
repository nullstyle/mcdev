package pkgindex

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// StdLibIndex indexes the golang stdlib packages
type StdLibIndex struct {
	Goroot string
	idx    *ManualIndex
}

// Index indexes all packages as returned by "go list std"
func (idx *StdLibIndex) Index() error {
	if idx.idx == nil {
		idx.idx = &ManualIndex{}
	}

	goroot := idx.Goroot
	if goroot == "" {
		goroot = os.Getenv("GOROOT")
	}

	gobin := filepath.Join(goroot, "bin", "go")

	out, err := exec.Command(gobin, "list", "std").Output()
	if err != nil {
		return err
	}

	packages := strings.Split(string(out), "\n")
	return idx.idx.Add(packages...)
}

//RefreshIfNeeded TODO
func (idx *StdLibIndex) RefreshIfNeeded() error {
	return nil
}

// Search returns any packages that match the query in the go standard library.
func (idx *StdLibIndex) Search(query string) (results []string, err error) {
	return idx.idx.Search(query)
}
