// Package cmd provides the common flags for all mcdev commands
package cmd

import (
	"flag"
)

// IsGB signifies that this command is being run within the root of a GB project
// each command should modify it'sgbehavior appropriately to support db
var IsGB = flag.Bool(
	"gb",
	false,
	"determine changed packages using the gb build tool's project layout",
)
