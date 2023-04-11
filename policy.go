package pagecache

import (
	"net/http"
	"time"

	"git.sr.ht/~jamesponddotco/pagecache-go/internal/httputil"
	"git.sr.ht/~jamesponddotco/pagecache-go/internal/sliceutil"
)

const (
	// DefaultTTL is the default time-to-live for cached responses.
	DefaultTTL time.Duration = 60 * time.Minute

	// DefaultMaxBodySize is the default maximum size of a response body.
	DefaultMaxBodySize int64 = 5 * 1024 * 1024
)

// Policy defines under which conditions an HTTP response may be cached.
type Policy struct {
	// AllowedStatusCodes is a list of HTTP status codes that should be cached.
	AllowedStatusCodes []int

	// AllowedMethods is a list of HTTP methods that should be cached.
	AllowedMethods []string

	// ExcludedHeaders is a list of HTTP headers to exclude from caching.
	ExcludedHeaders []string

	// ExcludedCookies is a list of HTTP cookies to exclude from caching.
	ExcludedCookies []string

	// Rules is a list of rules to apply to a request in order to determine if it should be cached.
	Rules []*Rule

	// MaxBodySize is the maximum size of the response body allowed to be
	// cached, in bytes. Zero or a negative value indicates no limit.
	MaxBodySize int64

	// UseCacheControl controls whether the cache takes the Cache-Control header
	// into account when deciding whether to cache a response.
	UseCacheControl bool

	// DefaultTTL is the default time-to-live of a cached response. Zero or a
	// negative value is interpreted as no expiration.
	//
	// If UseCacheControl is true, the cache will use the header's value to
	// determine the TTL instead.
	DefaultTTL time.Duration
}

// DefaultPolicy returns a new *Policy with opinionated but sane defaults.
func DefaultPolicy() *Policy {
	return &Policy{
		AllowedStatusCodes: []int{
			http.StatusOK,
			http.StatusNonAuthoritativeInfo,
			http.StatusNoContent,
			http.StatusPartialContent,
			http.StatusMultipleChoices,
			http.StatusMovedPermanently,
			http.StatusFound,
			http.StatusMethodNotAllowed,
			http.StatusGone,
			http.StatusRequestURITooLong,
			http.StatusNotImplemented,
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodHead,
		},
		ExcludedHeaders: []string{
			"Authorization",
		},
		ExcludedCookies: []string{
			"sessionid",
		},
		Rules:           []*Rule{},
		MaxBodySize:     DefaultMaxBodySize,
		DefaultTTL:      DefaultTTL,
		UseCacheControl: true,
	}
}

// IsCacheable checks if a given request and response pair is cacheable according
// to the policy. It evaluates status codes, methods, headers, cookies, and rules.
//
// Returns true if the request and response should be cached, otherwise false.
func (p *Policy) IsCacheable(resp *http.Response) bool {
	if !sliceutil.ContainsInt(p.AllowedStatusCodes, resp.StatusCode) {
		return false
	}

	if !sliceutil.ContainsString(p.AllowedMethods, resp.Request.Method) {
		return false
	}

	for header := range resp.Header {
		if sliceutil.ContainsString(p.ExcludedHeaders, header) {
			return false
		}
	}

	for _, cookie := range resp.Cookies() {
		if sliceutil.ContainsString(p.ExcludedCookies, cookie.Name) {
			return false
		}
	}

	if !httputil.IsBodySizeWithinLimit(resp.Header, p.MaxBodySize) {
		return false
	}

	// Check if the request URL matches any rules
	for _, rule := range p.Rules {
		if rule.Match(resp.Request.URL.String()) {
			if rule.Behavior == BehaviorExclude {
				return false
			}

			break
		}
	}

	// Check if Cache-Control should be taken into account
	if p.UseCacheControl {
		cacheControl := resp.Header.Get("Cache-Control")
		if cacheControl == "no-store" || cacheControl == "private" {
			return false
		}
	}

	// Check if the TTL is zero or negative
	ttl := p.TTL(resp)
	if ttl < 0 { //nolint:gosimple // explicitly checking for negative values to keep consistency
		return false
	}

	return true
}

// TTL returns the time-to-live (TTL) for the given response according to the
// policy. If the policy is configured to use the Cache-Control header and the
// header contains a valid max-age directive, the TTL will be based on that value.
// Otherwise, the policy's default TTL will be used.
func (p *Policy) TTL(resp *http.Response) time.Duration {
	if p.UseCacheControl {
		maxAge := httputil.MaxAge(resp.Header)

		if maxAge != -1 {
			return time.Duration(maxAge) * time.Second
		}
	}

	return p.DefaultTTL
}
