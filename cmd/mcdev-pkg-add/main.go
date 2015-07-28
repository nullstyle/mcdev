package main

// mcdev-pkg-add is a tool to help your go development workflow.  It:
//
// - ensures that the provided packages are imported into the provided file
//
// Usage:  `mcdev-pkg-add <file> [package ...]

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
	}

	file := args[0]
	log.Println(file)

	// parse the file,

	for _, pkg := range args[1:] {
		log.Println(pkg)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] path [import ...]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}
