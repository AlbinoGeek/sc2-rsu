// Code generated by go generate; DO NOT EDIT.
// This file was generated at: 2020-12-29 14:17:10.54064782 -0800 PST m=+0.000590627
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
const VERSION = "v0.3-94-gb4b93bf"

func init() {
	sc2replaystats.ClientIdentifier = fmt.Sprintf("%s-%s", PROGRAM, runtime.GOOS)
}
