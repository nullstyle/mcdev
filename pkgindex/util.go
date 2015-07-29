package pkgindex

import (
	"regexp"
	"strings"
)

const notEnd = "[^\x00]*?"

// searchToRegex converts the provided search string into a regular expression
// that can be used to match against a suffixarray in the manner required by this
// package. See documentation on the Searchable interface for details of the
// matching behavior provided by this package.
func searchToRegex(search string) (*regexp.Regexp, error) {
	letters := strings.Split(search, "")
	escaped := make([]string, len(letters))
	for i, l := range letters {
		escaped[i] = regexp.QuoteMeta(l)
	}

	expr := strings.Join([]string{
		"\x00",
		notEnd,
		strings.Join(escaped, notEnd),
		notEnd,
		"[^\x00]*",
	}, "")

	return regexp.Compile(expr)
}
