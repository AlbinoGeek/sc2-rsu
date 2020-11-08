package utils

import (
	"strconv"
	"strings"
	"time"
)

// ts2010 is a unix timestamp representing 2010-01-01 00:00:00
var ts2010 = time.Unix(1262304000, 0)

// ValidAPIKey checks that the given string matches sc2replaystats API key format
func ValidAPIKey(key string) bool {
	parts := strings.Split(key, ";")

	// there are exactly thee parts
	if len(parts) != 3 {
		return false
	}

	// the first two parts are each 40 characters
	if len(parts[0]) != 40 ||
		len(parts[1]) != 40 {
		return false
	}

	// the final part is a valid unix time
	ts, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return false
	}

	// time must be after 2010 (sanity check)
	if time.Unix(ts, 0).Before(ts2010) {
		return false
	}

	return true
}
