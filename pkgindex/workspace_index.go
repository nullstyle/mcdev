package pkgindex

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// WorkspaceIndex represents an index over a single go workspace, i.e. a
// directory with a `src` child directory that contains go source files arranged
// into packages.  This workspace would normally be an element of the curreny
// process's go path or the root of a gb project
type WorkspaceIndex struct {
	Dir string

	loadedAt time.Time
	idx      *ManualIndex
}

// Index indexes every package in underneath the `src` directory of Dir
func (idx *WorkspaceIndex) Index() error {
	if idx.idx == nil {
		idx.idx = &ManualIndex{}
	}

	srcDir := filepath.Join(idx.Dir, "src")
	packages := []string{}

	err := filepath.Walk(srcDir, func(path string, stat os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !stat.IsDir() {
			return nil
		}

		if filepath.Base(path)[0] == '.' {
			return filepath.SkipDir
		}

		hasGo, err := idx.hasGo(path)
		if err != nil {
			return err
		}

		if !hasGo {
			return nil
		}

		pkg := path[len(srcDir):]

		packages = append(packages, pkg)
		return nil
	})

	if err != nil {
		return err
	}

	return idx.idx.Add(packages...)
}

func (idx *WorkspaceIndex) RefreshIfNeeded() error {
	//TODO
	return nil
}

// hasGo returns true if the directory has any .go files
func (idx *WorkspaceIndex) hasGo(path string) (bool, error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, f := range entries {
		if f.IsDir() {
			continue
		}

		if filepath.Ext(f.Name()) == ".go" {
			return true, nil
		}
	}
	return false, nil
}
