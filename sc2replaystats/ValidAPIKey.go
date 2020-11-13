package sc2replaystats

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

	// there are exectly three parts, and the first two are 40 characters each
	if len(parts) != 3 ||
		len(parts[0]) != 40 ||
		len(parts[1]) != 40 {
		return false
	}

	// the final part is a valid unix time after 2010
	ts, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil ||
		time.Unix(ts, 0).Before(ts2010) {
		return false
	}

	return true
}
