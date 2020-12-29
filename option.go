package circleci

import (
	"net/http"
)

// Option is used to configure client with options
type Option func(c *Client)

// WithHTTPClient optionally sets the http.Client.
// This can be used when would like to customize HTTP option.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		if hc != nil {
			c.HTTPClient = hc
		}
	}
}

// WithPathPrefix optionally sets the API Prefix.
// This can be used when would like to mock or use different prefix.
func WithPathPrefix(path string) Option {
	return func(c *Client) {
		c.pathPrefix = path
	}
}
