package pkgindex

// Searchable types can filter their contained packages by the provided search
// string.  Any package that matches the search string (described below) should
// get returned.
//
// Matching is similar to that of a fuzzy-find in text editor:  an entry matches
// the search string if each letter in the search is present in the entry in the
// relative order specified by the string.  For example:
//
//   "aa" matches any package that has two a's in the import path
//   "ab" matches any package that has a "b" in its import path that follows any
//        "a" in its import path
//   "abc" matches any package that has an "a" that precedes a "b", which in turn
//        preceds an "c"
//
// These matches are performed using regular expressions derived from the search
// string.  The examples above would compile into the following regular
// expressions respectively:
//
//    ^.*a.*a.*$
//    ^.*a.*b.*$
//    ^.*a.*b.*c.*$
//
type Searchable interface {
	Search(string) ([]string, error)
}

// Index is the common interface to the various index structs
type Index interface {
	Searchable

	Index() error
	RefreshIfNeeded() error
}
