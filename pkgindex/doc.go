// Package pkgindex provides functions to index and query a local system's
// available go packages, allowing for quicker retrieval than traversing the
// file system.  Suitable for use in development tools
//
// This package has support for indexing normal go workspaces, but also provides
// a "gb" mode that is compatible with projects based upon the gb build tool.
package pkgindex
