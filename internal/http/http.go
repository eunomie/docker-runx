package http

import (
	"net/http"

	"github.com/eunomie/docker-runx/internal/constants"
)

type (
	userAgentTransporter struct {
		ua string
		rt http.RoundTripper
	}
)

func (u *userAgentTransporter) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", u.ua)
	return func(*http.Request) (*http.Response, error) {
		res, err := u.rt.RoundTrip(req)
		// Easy to put a breakpoint here to introspect HTTP headers on the response etc
		return res, err
	}(req)
}

func Transport() http.RoundTripper {
	return &userAgentTransporter{
		ua: constants.UserAgent,
		rt: http.DefaultTransport,
	}
}
