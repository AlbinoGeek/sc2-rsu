package utils

import "testing"

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
	for i, c := range cases {
		if res := CompareSemVer(c.Old, c.New); res != c.Expect {
			t.Errorf("Case %d failed: \"%s\" > \"%s\" = %d (expected: %d)", 1+i, c.New, c.Old, res, c.Expect)
		}
	}
}
