package sc2replaystats

import (
	"fmt"
	"runtime"
)

var (
	// Hostname represents the root domain sc2replaystats is hosted at
	Hostname = "sc2replaystats.com"

	// Protocol represents the HTTP protocol we use when communicating with sc2replaystats
	Protocol = "https"

	// APIRoot represents the base URL for requests to the sc2replaystats JSON-ish API
	APIRoot = fmt.Sprintf("%s://%s", Protocol, Hostname)

	// WebRoot represents the base URL for requests to the sc2replaystats Website
	WebRoot = fmt.Sprintf("%s://api.%s", Protocol, Hostname)

	// ClientIdentifier represents the "upload_method" shown to sc2replaystats
	ClientIdentifier = fmt.Sprintf("sc2-rsu-%s", runtime.GOOS)
)
