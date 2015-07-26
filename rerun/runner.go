package rerun

import (
	"log"
	"os/exec"
	"time"
)

// Runner attempts to keep the provided command running.  After Start() is
// called the underylying process will be restarted as requested as well as each
// time it exits.
//
// The configured cooldown will trigger when a restart is triggered due to the
// process exiting.  Manually triggered restarts--calling Restart()--will occur
// immediately.
type Runner struct {
	LastErr  error
	cmd      *exec.Cmd
	cooldown time.Duration

	exit     chan error
	restart  chan bool
	shutdown chan struct{}
	dontWait bool
	finished bool
	current  *exec.Cmd
}

// NewRunner constructs a new rerun service
func NewRunner(cmd *exec.Cmd, cooldown time.Duration) *Runner {
	return &Runner{
		cmd:      cmd,
		cooldown: cooldown,
		exit:     make(chan error, 1),
		restart:  make(chan bool, 1),
		shutdown: make(chan struct{}),
		current:  nil,
	}
}

// Start causes the underlying Cmd to be started
func (r *Runner) Start() {
	if r.current != nil {
		return
	}

	go r.run()
}

// Restart causes the undelying process to stop and restart immediately
func (r *Runner) Restart() {
	r.dontWait = true
	r.restart <- true
}

// Shutdown causes the undelying process to stop and restart immediately
func (r *Runner) Shutdown() {
	r.shutdown <- struct{}{}
}

func (r *Runner) run() {
	r.start()

	for {
		select {
		case err := <-r.exit:
			shouldRestart := r.finishProcess(err)
			r.current = nil

			if shouldRestart {
				r.start()
			} else {
				return
			}
		case <-r.restart:
			r.stop()
		case <-r.shutdown:
			log.Println("shutting down")
			r.finished = true
			r.stop()
		}
	}
}

func (r *Runner) finishProcess(err error) bool {
	var msg interface{}
	var isFatal bool

	if err == nil {
		msg = "exitted successfully"
		isFatal = false
	} else if _, ok := err.(*exec.ExitError); ok {
		msg = err
		isFatal = false
	} else {
		msg = err
		isFatal = true
	}

	if isFatal {
		log.Fatalln(msg)
	}

	log.Println(msg)

	return !r.finished
}

func (r *Runner) stop() {
	log.Println("stopping service")
	r.current.Process.Kill()
}

func (r *Runner) start() {
	next := *r.cmd
	r.current = &next

	if !r.dontWait {
		<-time.After(r.cooldown)
	}
	r.dontWait = false

	go func() {
		log.Println("starting service")
		r.exit <- r.current.Run()
	}()
}
