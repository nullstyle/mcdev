// Package pkgwatch provides the Watcher struct, which can watch a directory and
// emit events everytime the go code within the directory (or any sub-directory)
// changes.
//
// Useful for triggering commands that take a go package import path as an
// argument.
package pkgwatch
