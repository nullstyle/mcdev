package main

// mcdev-each-change is a tool to help your go development workflow.  It:
//
// - watches for .go files being changed underneath the current directory (recursively)
// - executes the templated command each time a package is changed
// - debounces executions by a configurable duration to allow for things like
//   gofmt to run prior to kicking the command off.  This is the `debounce` flag
// - provides a configurable cooldown for command executions to provide a
//   maximum rate of churn.
//
// This tool was designed to support a TDD-based development workflow that
// tests and re-installs a package everytime it is changed.  To do this, you would
// run:
//
// 		mcdev-each-change bash -c "go test {{.Pkg}} && go install {{.Pkg}}"
//
// The command will run until interupted using ctrl+c
//

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/nullstyle/mcdev/cmdtmpl"
	"github.com/nullstyle/mcdev/dotenv"
	"github.com/nullstyle/mcdev/pkgwatch"
	"github.com/nullstyle/mcdev/pkgwork"

	c "github.com/nullstyle/mcdev/cmd"
)

var done = make(chan os.Signal, 1)

var cmd *cmdtmpl.Command

var debounce = flag.Duration("debounce", 500*time.Millisecond, "how long to debounce package changes")
var cooldown = flag.Duration("cooldown", 4*time.Second, "how long to cooldown each command execution")

func main() {
	var err error

	flag.Parse()
	dotenv.Load()
	signal.Notify(done, os.Interrupt, os.Kill)

	cmd, err = cmdtmpl.NewCommand(flag.Args())
	if err != nil {
		log.Println("error when parsing command")
		log.Fatal(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Println("error when getting working directory")
		log.Fatal(err)
	}

	watcher := &pkgwatch.Watcher{
		Dir:      dir,
		Debounce: *debounce,
		IsGB:     *c.IsGB,
	}

	worker := &pkgwork.Worker{
		Fn:       execute,
		Cooldown: *cooldown,
	}
	if err := watcher.Run(); err != nil {
		log.Println("error when starting watcher")
		log.Fatal(err)
	}
	defer watcher.Close()

	log.Println("waiting for changes")

	for {
		select {
		case pkg := <-watcher.Changes():
			go func() {
				if err := worker.Run(pkg); err != nil {
					log.Fatal(err)
				}
			}()
		case _ = <-done:
			log.Println("shutting down")
			os.Exit(0)
		}
	}
}

func execute(pkg string) error {
	err := cmd.Run(struct{ Pkg string }{pkg})
	if err == nil {
		return nil
	}

	eerr, ok := err.(*exec.ExitError)
	if !ok {
		return err
	}

	log.Println(eerr)
	return nil
}
