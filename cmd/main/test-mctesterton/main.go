package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"text/template"
	"time"

	"github.com/nullstyle/testy-mctesterton/pkgwatch"
	"github.com/nullstyle/testy-mctesterton/pkgwork"
)

var done = make(chan os.Signal, 1)

var cmdStr = flag.String("cmd", "", "command line to execute upon package source change")
var cmdTmpl *template.Template

var debounce = flag.Duration("debounce", 500*time.Millisecond, "how long to debounce package changes")
var cooldown = flag.Duration("cooldown", 4*time.Second, "how long to cooldown each command execution")

func main() {
	var err error

	flag.Parse()
	signal.Notify(done, os.Interrupt, os.Kill)

	cmdTmpl, err = template.New("cmd").Parse(*cmdStr)
	if err != nil {
		log.Fatal(err)
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
	var cmdBuf bytes.Buffer
	err := cmdTmpl.Execute(&cmdBuf, struct{ Pkg string }{pkg})
	if err != nil {
		return err
	}

	full := cmdBuf.String()
	split := strings.Split(full, " ")
	cmd := exec.Command(split[0], split[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
