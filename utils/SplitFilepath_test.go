package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlbinoGeek/sc2-rsu/utils"
)

func TestSplitFilepath(t *testing.T) {
	var cases = []struct {
		Input string
		Ext   string
		Fname string
		Path  string
	}{
		{"/usr/local/src.a", ".a", "src", "/usr/local"},
		{"/.fstab", ".fstab", "", "/"},
		// TODO: Support Windows paths/roots
		// {"C:\\pagefile.sys", "sys", "pagefile", "C:\\"},
	}

	for _, c := range cases {
		p, f, e := utils.SplitFilepath(c.Input)
		assert.Equal(t, p, c.Path, "path must match")
		assert.Equal(t, f, c.Fname, "filename must match")
		assert.Equal(t, e, c.Ext, "extention must match")
	}
}
