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
	"strings"

	"github.com/nullstyle/mcdev/pkgindex"
	// "github.com/nullstyle/mcdev/cmd"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
	}

	var idx pkgindex.Index
	// load the index
	// if gb, create an index from the current gb project
	// else, create an index from the current GOPATH environment

	results, err := idx.Search(args[0])
	if err != nil {
		log.Fatal(err)
	}

	rs := strings.Join(results, "\n")
	fmt.Println(rs)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] query\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}
