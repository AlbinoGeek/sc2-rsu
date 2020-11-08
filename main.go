package main

import (
	"github.com/AlbinoGeek/sc2-rsu/cmd"
)

// build-time "constants" (Golang limitation: must be var of type string)
var (
	// PROGRAM is the human readable product name
	PROGRAM = "set-by-Makefile"

	// VERSION is the human readable product version
	VERSION = "set-by-Makefile"
)

func main() {
	cmd.PROGRAM = PROGRAM
	cmd.VERSION = VERSION
	cmd.Execute()
}
