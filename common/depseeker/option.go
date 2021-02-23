package depseeker

import "time"

// DefaultUserAgent default user agent
const DefaultUserAgent = "Depseeker - Open-source project (github.com/projectpandora/depseeker)"

// BlacklistPackageName blacklist
var BlacklistPackageName = []string{
	"version",
	"font-weight",
	"revision",
	"height",
	"z-index",
	"max-width",
	"ver",
	"site",
	"value",
	"latitude",
	"longitude",
	"weight",
	"height",
	"width",
	"ip",
	"id",
	"vtp_value",
	"namespace",
	"padding-top",
	"padding-bottom",
	"margin-top",
	"margin-bottom",
	"border-top-width",
	"border-bottom-width",
	"maximum-scale",
	"v",
	"website_id",
	"vtp_id",
}

// Options contains configuration options for the client
type Options struct {
	UserAgent  string
	HTTPProxy  string
	SocksProxy string
	Threads    int
	CdnCheck   bool
	// Timeout is the maximum time to wait for the request
	Timeout time.Duration
	// RetryMax is the maximum number of retries
	RetryMax            int
	CustomHeaders       map[string]string
	FollowRedirects     bool
	FollowHostRedirects bool
}

// DefaultOptions contains the default options
var DefaultOptions = Options{
	Threads:   25,
	Timeout:   30 * time.Second,
	RetryMax:  5,
	CdnCheck:  true,
	UserAgent: DefaultUserAgent,
}
