package pkgwatch

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// This file contains functions that help interface with a process's GOPATH
// environment variable

var gopath []string

func init() {
	raw := os.Getenv("GOPATH")
	unexpanded := strings.Split(raw, ":")

	for _, up := range unexpanded {
		ap, err := filepath.Abs(up)
		if err != nil {
			log.Fatal(err)
		}

		gopath = append(gopath, ap)
	}
}

// isOnGoPath searches the gopath, returning the first entry that is a parent of the
// provided absolute directory path.
func isOnGoPath(dir string) (string, bool) {
	for _, gp := range gopath {
		if strings.HasPrefix(dir, gp) {
			return gp, true
		}
	}

	return "", false
}
