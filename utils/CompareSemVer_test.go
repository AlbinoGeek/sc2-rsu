package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlbinoGeek/sc2-rsu/utils"
)

func TestCompareSemVer(t *testing.T) {
	var cases = []struct {
		Expect int
		New    string
		Old    string
	}{
		{-1, "", "1"},
		{1, "1", "0"},
		{1, "v1", "v0"},
		{1, "1", ""},
		{-1, "v-beta", "1"},
		{1, "1.2", "1-alpha"},
		{-1, "1-beta", "1.1.2"},
	}

	for _, c := range cases {
		assert.Equal(t, utils.CompareSemVer(c.Old, c.New), c.Expect, "result must match")
	}
}
