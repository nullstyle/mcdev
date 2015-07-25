package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/nullstyle/testy-mctesterton/pkgwatch"
)

var done = make(chan os.Signal, 1)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	watcher, err := pkgwatch.NewWatcher(dir)
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	signal.Notify(done, os.Interrupt, os.Kill)
	watcher.Run()

	log.Println("waiting for changes")

	for {
		select {
		case pkg := <-watcher.Changes():
			log.Printf("changed: %s", pkg)
		case _ = <-done:
			log.Println("shutting down")
			os.Exit(0)
		}
	}
}
