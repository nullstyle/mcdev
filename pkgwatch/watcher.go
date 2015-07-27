package pkgwatch

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

var isGB = flag.Bool("gb", false, "determine changed packages using the gb build tool's project layout")

// Watcher watches for go package changes underneath a directory and emits
// their names as go files within them change
type Watcher struct {
	Dir      string
	Debounce time.Duration

	inited  bool
	fs      *fsnotify.Watcher
	changes chan string
	done    chan bool
	pending map[string]bool
}

// Init ensures the internal state of the watcher is properly initialized
func (w *Watcher) Init() (err error) {
	if w.inited {
		return nil
	}

	w.fs, err = fsnotify.NewWatcher()
	if err != nil {
		return
	}

	stat, err := os.Stat(w.Dir)
	if err != nil {
		return
	}

	if !stat.IsDir() {
		err = fmt.Errorf("%s is not a directory", w.Dir)
		return
	}

	w.changes = make(chan string, 100)
	w.done = make(chan bool, 1)
	w.pending = make(map[string]bool)
	w.inited = true
	return
}

// Run runs the watcher, continually pushing events from the fs watcher to
// the changes channel.
func (w *Watcher) Run() error {
	if err := w.Init(); err != nil {
		return err
	}

	// initialize the watchlist
	err := filepath.Walk(w.Dir, func(path string, stat os.FileInfo, err error) error {

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
			case <-time.After(w.Debounce):
				w.emit()
			case <-w.done:
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
	base := filepath.Base(path)
	// we don't watch hidden dirs
	if base[0] == '.' {
		return filepath.SkipDir
	}

	// we don't watch vendored code for modifications (TODO: add a flag to enable/disable this)
	if base == "vendor" {
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

	w.addPending(pkg)
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
	var foundRoot string

	if *isGB {
		foundRoot = w.Dir
	} else {
		var foundOnGoPath bool
		foundRoot, foundOnGoPath = isOnGoPath(dir)

		if !foundOnGoPath {
			log.Printf("warn: changed file was not found on a local gopath. ignoring...")
			return "", false
		}
	}

	srcRoot := filepath.Join(foundRoot, "src")

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

func (w *Watcher) emit() {
	if len(w.pending) == 0 {
		return
	}

	for pkg := range w.pending {
		w.changes <- pkg
	}
	w.pending = make(map[string]bool)
}

func (w *Watcher) addPending(pkg string) {
	w.pending[pkg] = true
}
