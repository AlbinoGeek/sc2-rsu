// Code generated by go generate; DO NOT EDIT.
// This file was generated at: 2020-12-29 14:18:23.167253519 -0800 PST m=+0.000461375
// This file is in sync with: 030a488a9e69ea1c0a23848f6a1a3d3bdde95cb6
package cmd

import (
	"fmt"
	"runtime"

	"github.com/AlbinoGeek/sc2-rsu/sc2replaystats"
)

// PROGRAM is the human readable product name
const PROGRAM = "sc2-rsu"

// VERSION is the human readable product version
const VERSION = "v0.4"

func init() {
	sc2replaystats.ClientIdentifier = fmt.Sprintf("%s-%s", PROGRAM, runtime.GOOS)
}
