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
	started map[string]time.Time
	running map[string]bool
	again   map[string]bool
}

func (w *Worker) Init() {
	if w.inited {
		return
	}

	w.started = map[string]time.Time{}
	w.running = map[string]bool{}
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
	again := w.again[pkg]
	w.again[pkg] = false
	w.Unlock()
	return again
}

func (w *Worker) shouldStart(pkg string) bool {
	w.Lock()

	startedAt := w.started[pkg]
	alreadyRunning := w.running[pkg]
	coolEnough := time.Since(startedAt) > w.Cooldown

	if !coolEnough {
		w.Unlock()
		return false
	}

	if alreadyRunning {
		log.Printf("requeue: %s", pkg)
		w.again[pkg] = true

		w.Unlock()
		return false
	}

	w.started[pkg] = time.Now()
	w.running[pkg] = true
	w.Unlock()
	return true
}
