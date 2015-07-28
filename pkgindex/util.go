package pkgindex

import (
	"log"
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
	expr := strings.Join([]string{
		"\x00",
		notEnd,
		strings.Join(letters, notEnd),
		notEnd,
		"[^\x00]*",
	}, "")

	log.Println(expr)

	return regexp.Compile(expr)
}
