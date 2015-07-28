package pkgindex

import (
	"index/suffixarray"
	"strings"
)

// ManualIndex represents an index of manually added packages.  Calling "Add"
type ManualIndex struct {
	packages []string
	index    *suffixarray.Index
}

// Add appends the provided packages onto the index and triggers a re-index.
func (idx *ManualIndex) Add(packages ...string) error {
	idx.packages = append(idx.packages, packages...)
	return idx.Index()
}

// Index builds a new suffixarray for the package names previously registered
// with this instance.
func (idx *ManualIndex) Index() error {
	var input = "\x00" + strings.Join(idx.packages, "\x00")
	idx.index = suffixarray.New([]byte(input))
	return nil
}

// RefreshIfNeeded is a no-op for ManualIndex instances.  They are always up to
// date.
func (idx *ManualIndex) RefreshIfNeeded() error {
	return nil
}

// Search returns any packages that match the query.  See documentation for the
// Searchable interface for details.
func (idx *ManualIndex) Search(query string) (results []string, err error) {
	regexp, err := searchToRegex(query)
	if err != nil {
		return nil, err
	}

	positions := idx.index.FindAllIndex(regexp, -1)
	bytes := idx.index.Bytes()

	for _, pos := range positions {
		start := pos[0] + 1 //we add one to strip of the leading null byte
		end := pos[1]
		results = append(results, string(bytes[start:end]))
	}
	return
}
