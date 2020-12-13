package utils_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlbinoGeek/sc2-rsu/utils"
)

func TestDownloadWriter(t *testing.T) {
	var cases = []struct {
		URL    string
		Result string
		ResErr bool
	}{
		{"/usr/local/src.a", "", true}, // Not a URL
		{"http://non-exist-url.loltld", "", true},
		{"http://example.org/non-exist", "", true}, //  status != 200
		{"https://raw.githubusercontent.com/AlbinoGeek/sc2-rsu/main/_dist/.gitkeep", "", false},
	}

	for _, c := range cases {
		w := &bytes.Buffer{}

		err := utils.DownloadWriter(c.URL, w)

		// ! this is wrong I'm sure
		if c.ResErr {
			assert.NotEqual(t, err, nil, "must error")
		} else {
			assert.Equal(t, err, nil, "must not error")
		}

		assert.Equal(t, w.String(), c.Result, "result must match")
	}
}
