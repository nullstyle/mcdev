package pkgwork

import (
	"log"
	"sync"
	"time"
)

type Worker struct {
	Fn       func(string) error
	Cooldown time.Duration

	sync.Mutex

	inited  bool
	running map[string]time.Time
	again   map[string]bool
}

func (w *Worker) Init() {
	if w.inited {
		return
	}

	w.running = map[string]time.Time{}
	w.again = map[string]bool{}
	w.inited = true
}

func (w *Worker) Run(pkg string) error {
	w.Init()

	if !w.shouldStart(pkg) {
		return nil
	}

	defer w.finish(pkg)

	err := w.run(pkg)
	if err != nil {
		return err
	}

	for w.shouldRunAgain(pkg) {
		err := w.run(pkg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Worker) run(pkg string) error {
	log.Printf("starting: %s", pkg)
	defer log.Printf("done: %s", pkg)
	return w.Fn(pkg)
}

func (w *Worker) finish(pkg string) {
	w.Lock()
	delete(w.running, pkg)
	delete(w.again, pkg)
	w.Unlock()
}

func (w *Worker) shouldRunAgain(pkg string) bool {
	w.Lock()
	again, _ := w.again[pkg]
	w.again[pkg] = false
	w.Unlock()
	return again
}

func (w *Worker) shouldStart(pkg string) bool {
	w.Lock()

	startedAt, alreadyRunning := w.running[pkg]
	if alreadyRunning {

		// If the cooldown has elapsed since the last start of the work for this
		// pkg, requeue the pkg to run again after the current process is complete
		if time.Since(startedAt) > w.Cooldown {
			log.Printf("requeuing %s", pkg)
			w.again[pkg] = true
		}

		w.Unlock()
		return false
	}

	w.running[pkg] = time.Now()
	w.Unlock()
	return true
}
