// Package pkgindex provides functions to index and query a local system's
// available go packages, allowing for quicker retrieval than traversing the
// file system.  Suitable for use in development tools
//
// This package has support for indexing normal go workspaces, but also provides
// a "gb" mode that is compatible with projects based upon the gb build tool.
//
// This package contains a set of structures that represent various types of
// indexes.  Notably there is the GBIndex struct and the GoPathIndex struct which
// wrap the other index types to provide an index over a gb project or the local
// system's GOPATH, respectively.
package pkgindex
