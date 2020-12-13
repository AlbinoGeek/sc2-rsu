package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlbinoGeek/sc2-rsu/utils"
)

func TestStripPathParts(t *testing.T) {
	var cases = []struct {
		Path   string
		Strip  int
		Result string
	}{
		{"/home/foo", 1, "/home"},
		{"/home/foo", 2, "/"},
		{"/foo", -1, "foo"},                  // 1: the root is a path part!
		{"/foo", 1, "/"},                     // (1)
		{"/foo", -3, ""},                     // 2: more path parts than we have
		{"/foo", 3, "/"},                     // (1 & 2)
		{"/home/foo/", 1, "/home"},           // 3: strip last character if it's a separator
		{"/usr/var/foo/", -1, "usr/var/foo"}, // 4: stripLeft creates relative paths
		{"/usr/var/foo/", -2, "var/foo"},     // (4)
		{"/usr/var/foo/", -3, "foo"},         // (4)
		{"usr/var", 1, "usr"},                // 5: relative paths stay relative
		{"usr/var", 2, ""},                   // (5)
	}

	for _, c := range cases {
		assert.Equal(t, utils.StripPathParts(c.Path, c.Strip), c.Result, "result must match")
	}
}
