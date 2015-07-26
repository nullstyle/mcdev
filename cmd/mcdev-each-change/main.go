package main

// mcdev-each-change is a tool to help your go development workflow.  It:
//
// - watches for .go files being changed underneath the current directory (recursively)
// - executes the templated command (the `cmd` flag) each time a package is changed
// - debounces executions by a configurable duration to allow for things like
//   gofmt to run prior to kicking the command off.  This is the `debounce` flag
// - provides a configurable cooldown for command executions to provide a
//   maximum rate of churn.
//
// This tool was designed to support a TDD-based development workflow that
// tests and re-installs a package everytime it is changed.  To do this, you would
// run:
//
// 		mcdev-each-change -cmd="go test {{.Pkg}} && go install {{.Pkg}}"
//
// The command will run until interupted using ctrl+c
//

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"text/template"
	"time"

	"github.com/nullstyle/testy-mctesterton/pkgwatch"
	"github.com/nullstyle/testy-mctesterton/pkgwork"
)

var done = make(chan os.Signal, 1)

var cmdTmpls []*template.Template

var debounce = flag.Duration("debounce", 500*time.Millisecond, "how long to debounce package changes")
var cooldown = flag.Duration("cooldown", 4*time.Second, "how long to cooldown each command execution")

func main() {
	var err error

	flag.Parse()
	signal.Notify(done, os.Interrupt, os.Kill)

	for _, arg := range flag.Args() {
		t, err := template.New("cmd").Parse(arg)
		if err != nil {
			log.Fatal(err)
		}
		cmdTmpls = append(cmdTmpls, t)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	watcher := &pkgwatch.Watcher{
		Dir:      dir,
		Debounce: *debounce,
	}

	worker := &pkgwork.Worker{
		Fn:       execute,
		Cooldown: *cooldown,
	}

	if err := watcher.Run(); err != nil {
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

	args := make([]string, len(cmdTmpls))
	for i, t := range cmdTmpls {
		var buf bytes.Buffer

		err := t.Execute(&buf, struct{ Pkg string }{pkg})
		if err != nil {
			return err
		}
		args[i] = buf.String()
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
