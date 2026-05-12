// Package httpmirror provides an HTTP-backed Simple Streams source.
package httpmirror

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	simplestreams "github.com/meigma/go-simplestreams"
)

// Source opens Simple Streams content from an HTTP mirror root.
type Source struct {
	baseURL   *url.URL
	client    *http.Client
	userAgent string
}

// Option configures a Source.
type Option func(*Source)

// WithHTTPClient sets the HTTP client used by the source.
func WithHTTPClient(client *http.Client) Option {
	return func(source *Source) {
		if client != nil {
			source.client = client
		}
	}
}

// WithUserAgent sets the User-Agent header sent on requests.
func WithUserAgent(userAgent string) Option {
	return func(source *Source) {
		source.userAgent = userAgent
	}
}

// New constructs an HTTP source rooted at baseURL.
func New(baseURL string, options ...Option) (*Source, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, errors.New("httpmirror: base URL scheme must be http or https")
	}
	source := &Source{
		baseURL: parsed,
		client:  http.DefaultClient,
	}
	for _, option := range options {
		if option != nil {
			option(source)
		}
	}
	return source, nil
}

// Open issues a GET request for path under the source base URL.
func (source *Source) Open(ctx context.Context, path simplestreams.RelativePath) (io.ReadCloser, error) {
	if err := path.Validate(); err != nil {
		return nil, err
	}
	requestURL, err := url.JoinPath(source.baseURL.String(), path.String())
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	if source.userAgent != "" {
		request.Header.Set("User-Agent", source.userAgent)
	}

	response, err := source.client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		if response.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%w: GET %s: %s", simplestreams.ErrNotFound, requestURL, response.Status)
		}
		return nil, fmt.Errorf("httpmirror: GET %s: %s", requestURL, response.Status)
	}
	return response.Body, nil
}
