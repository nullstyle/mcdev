package pkgindex

import (
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

// indexes every package in underneath the `src` directory of Dir
func (idx *WorkspaceIndex) Index() error {
	srcDir := filepath.Join(idx.Dir, "src")

	//TODO
	_ = srcDir

	return nil
}

func (idx *WorkspaceIndex) RefreshIfNeeded() error {
	//TODO
	return nil
}
