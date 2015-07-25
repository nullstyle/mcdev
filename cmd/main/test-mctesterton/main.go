package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nullstyle/testy-mctesterton/pkgwatch"
)

var done = make(chan os.Signal, 1)

func main() {
	signal.Notify(done, os.Interrupt, os.Kill)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	watcher := &pkgwatch.Watcher{
		Dir:      dir,
		Debounce: 300 * time.Millisecond,
	}

	if err := watcher.Run(); err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

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
