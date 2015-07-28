package main

// mcdev-rerun is a tool to help your go development workflow.  It:
//
// - starts the command provided
// - watches for .go files being changed underneath the current directory (recursively)
// - debounces restarts by a configurable duration to allow for things like
//   gofmt to run prior to restarting the service.  This is the `debounce` flag
// - restarts the command provided anytime it exits
//
// This tool was designed to support a development workflow for a server process
// where you would like to restart the server or other-long running process
// anytime one of it's package dependencies change.
//
// run:
//
// 		mcdev-rerun go run myserver.go
//
// The server process will run until stopped using ctrl+c
//

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/nullstyle/mcdev/cmdtmpl"
	"github.com/nullstyle/mcdev/dotenv"
	"github.com/nullstyle/mcdev/pkgwatch"
	"github.com/nullstyle/mcdev/rerun"
)

var debounce = flag.Duration("debounce", 1*time.Second, "how long to debounce package changes")
var cooldown = flag.Duration("cooldown", 1*time.Second, "how long to cooldown each command execution")

var sigs = make(chan os.Signal, 1)
var lock sync.Mutex
var proc *rerun.Runner
var cmd *cmdtmpl.Command

func main() {
	var err error

	flag.Parse()
	dotenv.Load()
	signal.Notify(sigs, os.Interrupt, os.Kill)

	cmd, err = cmdtmpl.NewCommand(flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	p, err := cmd.Make(struct{}{})
	if err != nil {
		log.Fatal(err)
	}

	proc = rerun.NewRunner(p, *cooldown)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	watcher := &pkgwatch.Watcher{
		Dir:      dir,
		Debounce: *debounce,
	}

	if err := watcher.Run(); err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	proc.Start()

	for {
		select {
		case <-watcher.Changes():
			proc.Restart()
		case <-sigs:
			proc.Shutdown()
			os.Exit(0)
		}
	}
}
