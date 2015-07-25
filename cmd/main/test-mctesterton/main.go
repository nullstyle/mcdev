package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nullstyle/testy-mctesterton/pkgwatch"
	"github.com/nullstyle/testy-mctesterton/pkgwork"
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

	worker := &pkgwork.Worker{
		Fn: func(pkg string) error {
			log.Printf("working on %s", pkg)
			<-time.After(30 * time.Second)
			log.Printf("finished with %s", pkg)
			return nil
		},
		Cooldown: 5 * time.Second,
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
