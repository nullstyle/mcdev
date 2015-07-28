package pkgindex

// CompoundIndex searches each of its child indexes, merging the results
type CompoundIndex struct {
	children []Index
}

// Index calls through to all child indexes
func (idx *CompoundIndex) Index() error {
	return idx.each(func(cidx Index) error {
		return cidx.Index()
	})
}

// RefreshIfNeeded calls through to all child indexes.
func (idx *CompoundIndex) RefreshIfNeeded() error {
	return idx.each(func(cidx Index) error {
		return cidx.RefreshIfNeeded()
	})
}

// Search returns any packages that match the query in all child indexes.
func (idx *CompoundIndex) Search(query string) (results []string, err error) {
	alreadyAdded := map[string]bool{}

	err = idx.each(func(cidx Index) error {
		r, err := cidx.Search(query)
		if err != nil {
			return err
		}

		// append results to final result only in the case it hasn't been added previously
		for _, pkg := range r {
			if alreadyAdded[pkg] {
				continue
			}

			alreadyAdded[pkg] = true
			results = append(results, pkg)
		}
		return nil
	})
	return
}

// Add registers a child index, adding it to the set of indexes that will be
// queried when performing operations on this compound index.
func (idx *CompoundIndex) Add(cidx Index) {
	idx.children = append(idx.children, cidx)
}

// each performs the provided function on every child index, in order, until an
// error occurs.
func (idx *CompoundIndex) each(fn func(Index) error) error {
	for _, cidx := range idx.children {
		err := fn(cidx)
		if err != nil {
			return err
		}
	}
	return nil
}
