// Package staticresponse a demo plugin.
package staticresponse

import (
	"context"
	"fmt"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	StatusCode int         `json:"statusCode,omitempty"`
	Body       string      `json:"body,omitempty"`
	Headers    http.Header `json:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		StatusCode: http.StatusOK,
		Body:       "",
		Headers:    http.Header{},
	}
}

// StaticResponse plugin.
type StaticResponse struct {
	next       http.Handler
	statusCode int
	body       string
	headers    http.Header
	name       string
}

// New created a new StaticResponse plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.StatusCode < 100 || config.StatusCode > 999 {
		return nil, fmt.Errorf("invalid response status code")
	}

	// Ensure headers are not nil if left unconfigured
	if config.Headers == nil {
		config.Headers = http.Header{}
	}

	return &StaticResponse{
		next:       next,
		name:       name,
		statusCode: config.StatusCode,
		body:       config.Body,
		headers:    config.Headers,
	}, nil
}

// ServeHTTP function required to make StaticResponse comply with http.Handler interface.
func (s *StaticResponse) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Set headers first before sending the response
	for key, values := range s.headers {
		for _, value := range values {
			rw.Header().Add(key, value)
		}
	}

	rw.WriteHeader(s.statusCode)

	if s.body != "" {
		_, _ = rw.Write([]byte(s.body))
	}
}
