package rerun

import (
	"log"
	"os"
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

	exit      chan error
	restart   chan bool
	dontWait  bool
	finishing bool
	finished  bool
	current   *exec.Cmd
}

// NewRunner constructs a new rerun service
func NewRunner(cmd *exec.Cmd, cooldown time.Duration) *Runner {
	return &Runner{
		cmd:      cmd,
		cooldown: cooldown,
		exit:     make(chan error, 1),
		restart:  make(chan bool),
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

// Shutdown kills this service runner and waits until it is complete
func (r *Runner) Shutdown() {
	close(r.restart)

	done := make(chan bool, 1)
	go func() {
		for {
			if r.finished {
				done <- true
				return
			}
			<-time.After(1 * time.Second)
		}
	}()
	<-done
}

func (r *Runner) run() {
	r.start()

	for {
		select {
		case err := <-r.exit:
			r.finishProcess(err)
			r.current = nil
			r.start()
		case _, more := <-r.restart:
			// if the restart channel closed, shutdown
			if !more {
				r.shuttingDown()
				return
			}

			r.stop()
		}
	}
}

func (r *Runner) shuttingDown() {
	r.finishing = true
	r.stop()
	select {
	case err := <-r.exit:
		r.finishProcess(err)
		r.finished = true
		log.Println("shutdown complete")
	case <-time.After(10 * time.Second):
		log.Fatalln("shutdown did not complete before 10 seconds elapsed")
	}
}

func (r *Runner) finishProcess(err error) {
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
}

func (r *Runner) stop() {
	log.Println("stopping service")
	r.current.Process.Signal(os.Interrupt)
}

func (r *Runner) start() {
	if r.finishing {
		return
	}

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
