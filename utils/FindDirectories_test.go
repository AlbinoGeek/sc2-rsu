package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlbinoGeek/sc2-rsu/utils"
)

func TestFindDirectories(t *testing.T) {
	var cases = []struct {
		Root   string
		Suffix string
		Result []string
		ResErr bool
	}{
		{"/non-exist", "", nil, true},
		{"/root", "", []string{}, false},
	}

	// TODO: add test for Ignore
	for _, c := range cases {
		res, err := utils.FindDirectoriesBySuffix(c.Root, c.Suffix, true)

		// ! this is wrong I'm sure
		if c.ResErr {
			assert.NotEqual(t, err, nil, "must error")
		} else {
			assert.Equal(t, err, nil, "must not error")
		}

		l1 := len(res)
		l2 := len(c.Result)
		if l1 != l2 {
			t.Errorf("result mismatch (len=%d) != (len=%d)", l1, l2)
		}

		for i := range res {
			assert.Equal(t, res[i], c.Result[i], "results must match")
		}
	}

	// TODO: test for actual results...
}
