package pkgwatch

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-fsnotify/fsnotify"
)

// Watcher watches for go package changes underneath a directory and emits
// their names as go files within them change
type Watcher struct {
	dir     string
	fs      *fsnotify.Watcher
	changes chan string
	done    chan bool
	dirs    map[string]bool
}

// NewWatcher created a new package watcher under the provided dir
func NewWatcher(dir string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(dir)

	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}

	return &Watcher{
		dir,
		watcher,
		make(chan string, 10),
		make(chan bool, 1),
		map[string]bool{},
	}, nil
}

// Run runs the watcher, continually pushing events from the fs watcher to
// the changes channel.
func (w *Watcher) Run() error {
	// initialize the watchlist
	err := filepath.Walk(w.dir, func(path string, stat os.FileInfo, err error) error {

		if err != nil {
			log.Fatal(err)
		}

		if !stat.IsDir() {
			return nil
		}

		return w.AddPath(path)
	})

	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-w.fs.Events:
				// if it's a directory add/remove it from the watchlist
				err = w.processDirEvent(event)
				if err != nil {
					log.Fatal(err)
				}

				// if it's a go file, find what package
				err := w.processGoEvent(event)
				if err != nil {
					log.Fatal(err)
				}

			case err := <-w.fs.Errors:
				if err != nil {
					log.Fatal(err)
				}
			case _ = <-w.done:
				close(w.changes)
				log.Println("closed package watcher")
				return
			}
		}
	}()

	return nil
}

// Close closes the watcher
func (w *Watcher) Close() error {
	w.done <- true
	close(w.done)
	return w.fs.Close()
}

// Changes return a channel a message everytime a package underneath the
// watched directory changes
func (w *Watcher) Changes() <-chan string {
	return w.changes
}

// AddPath adds a new directory to the watched list of directories
func (w *Watcher) AddPath(path string) error {
	// we don't watch vendored code for modifications (TODO: add a flag to enable/disable this)
	if filepath.Base(path) == "vendor" {
		return filepath.SkipDir
	}

	if err := w.fs.Add(path); err != nil {
		return err
	}

	return nil
}

//processGoEvent takes a go package, finds out the its path relative to
//the current GOPATH, and emits it on the changes channel
func (w *Watcher) processGoEvent(event fsnotify.Event) error {
	goPath := event.Name

	if filepath.Ext(goPath) != ".go" {
		return nil
	}

	dir := filepath.Dir(goPath)
	pkg, found := w.findPackage(dir)

	if !found {
		log.Printf("couldn't find package for %s", filepath.Base(goPath))
		return nil
	}

	w.changes <- pkg
	return nil
}

func (w *Watcher) processDirEvent(event fsnotify.Event) error {
	if event.Op&fsnotify.Create != fsnotify.Create {
		return nil
	}

	stat, err := os.Stat(event.Name)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return nil
	}

	return w.AddPath(event.Name)
}

func (w *Watcher) findPackage(dir string) (string, bool) {
	gopath, found := isOnGoPath(dir)

	if !found {
		log.Printf("warn: changed file was not found on a local gopath. ignoring...")
		return "", false
	}

	srcRoot := filepath.Join(gopath, "src")

	// ASSERT: the changed directory is underneath the found gopath entry's "src"
	// directory
	if !strings.HasPrefix(dir, srcRoot) {
		log.Fatalf("%s isn't underneath %s, "+
			"even though it was found with isOnGoPath",
			dir,
			srcRoot)
	}

	return dir[len(srcRoot)+1:], true
}
