package sc2replaystats

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrBadKey is the base message shown in all errors returned by ValidAPIKey
	ErrBadKey = errors.New("invalid sc2replaystats API key")

	// ErrBadKeyDate means the API key's timestamp was invalid or too old to use
	ErrBadKeyDate = fmt.Errorf("%s: bad timestamp", ErrBadKey.Error())

	// ErrBadKeyLen means the API key's format or length was incorrect
	ErrBadKeyLen = fmt.Errorf("%s: bad length", ErrBadKey.Error())

	// ts2010 is a unix timestamp representing 2010-01-01 00:00:00
	ts2010 = time.Unix(1262304000, 0)
)

// ValidAPIKey checks that the given string matches sc2replaystats API key format
func ValidAPIKey(key string) bool {
	parts := strings.Split(key, ";")

	// there are exectly three parts, and the first two are 40 characters each
	if len(parts) != 3 ||
		len(parts[0]) != 40 ||
		len(parts[1]) != 40 {
		return false // ErrBadKeyLen
	}

	// the final part is a valid unix time after 2010
	ts, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil ||
		time.Unix(ts, 0).Before(ts2010) {
		return false // ErrBadKeyDate
	}

	return true
}
