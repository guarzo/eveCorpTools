package utils

import (
	"net/http"
)

// NewHTTPClientWithUserAgent creates an HTTP client that automatically includes the specified User-Agent header.
func NewHTTPClientWithUserAgent(userAgent string) *http.Client {
	// Define a custom RoundTripper that adds the User-Agent header
	customTransport := &userAgentRoundTripper{
		Wrapped:   http.DefaultTransport,
		UserAgent: userAgent,
	}

	return &http.Client{
		Transport: customTransport,
	}
}

// userAgentRoundTripper is a custom RoundTripper that adds a User-Agent header.
type userAgentRoundTripper struct {
	Wrapped   http.RoundTripper
	UserAgent string
}

// RoundTrip implements the RoundTripper interface.
func (rt *userAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	reqClone := req.Clone(req.Context())

	// Set the User-Agent header
	reqClone.Header.Set("User-Agent", rt.UserAgent)

	// Perform the request using the wrapped RoundTripper
	return rt.Wrapped.RoundTrip(reqClone)
}
