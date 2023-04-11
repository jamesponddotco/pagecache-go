package pagecache_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"git.sr.ht/~jamesponddotco/pagecache-go"
)

func TestPolicy_IsCacheable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		policy         *pagecache.Policy
		response       *http.Response
		expectedResult bool
	}{
		{
			name:   "IsCacheable with default policy",
			policy: pagecache.DefaultPolicy(),
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
				},
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Length": []string{"1000"},
				},
			},
			expectedResult: true,
		},
		{
			name:   "IsCacheable with status code not allowed",
			policy: pagecache.DefaultPolicy(),
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
				},
				StatusCode: http.StatusForbidden,
				Header: http.Header{
					"Content-Length": []string{"1000"},
				},
			},
			expectedResult: false,
		},
		{
			name:   "IsCacheable with method not allowed",
			policy: pagecache.DefaultPolicy(),
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodPost,
				},
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Length": []string{"1000"},
				},
			},
			expectedResult: false,
		},
		{
			name: "IsCacheable with excluded header",
			policy: func() *pagecache.Policy {
				p := pagecache.DefaultPolicy()
				p.ExcludedHeaders = append(p.ExcludedHeaders, "X-Custom")
				return p
			}(),
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
				},
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Length": []string{"1000"},
					"X-Custom":       []string{"test"},
				},
			},
			expectedResult: false,
		},
		{
			name: "IsCacheable with excluded cookie",
			policy: func() *pagecache.Policy {
				p := pagecache.DefaultPolicy()
				p.ExcludedCookies = append(p.ExcludedCookies, "testcookie")
				return p
			}(),
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
				},
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Length": []string{"1000"},
					"Set-Cookie":     []string{"testcookie=value"},
				},
			},
			expectedResult: false,
		},
		{
			name: "IsCacheable with max body size exceeded",
			policy: func() *pagecache.Policy {
				p := pagecache.DefaultPolicy()
				p.MaxBodySize = 500
				return p
			}(),
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
				},
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Length": []string{"1000"},
				},
			},
			expectedResult: false,
		},
		{
			name:   "IsCacheable with Cache-Control no-store",
			policy: pagecache.DefaultPolicy(),
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
				},
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Length": []string{"1000"},
					"Cache-Control":  []string{"no-store"},
				},
			},
			expectedResult: false,
		},
		{
			name:   "IsCacheable with Cache-Control private",
			policy: pagecache.DefaultPolicy(),
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
				},
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Length": []string{"1000"},
					"Cache-Control":  []string{"private"},
				},
			},
			expectedResult: false,
		},
		{
			name: "IsCacheable with rule exclude",
			policy: func() *pagecache.Policy {
				p := pagecache.DefaultPolicy()
				p.Rules = append(p.Rules, &pagecache.Rule{
					URL:      "http://example.com/excluded",
					Behavior: pagecache.BehaviorExclude,
				})
				return p
			}(),
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    parseTestURL(t, "http://example.com/excluded"),
				},
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Length": []string{"1000"},
				},
			},
			expectedResult: false,
		},
		{
			name: "IsCacheable with rule include",
			policy: func() *pagecache.Policy {
				p := pagecache.DefaultPolicy()
				p.Rules = append(p.Rules, &pagecache.Rule{
					URL:      "http://example.com/included",
					Behavior: pagecache.BehaviorInclude,
				})
				return p
			}(),
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    parseTestURL(t, "http://example.com/included"),
				},
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Length": []string{"1000"},
				},
			},
			expectedResult: true,
		},
		{
			name: "IsCacheable with negative TTL",
			policy: func() *pagecache.Policy {
				p := pagecache.DefaultPolicy()
				p.UseCacheControl = false
				p.DefaultTTL = -1 * time.Second
				return p
			}(),
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
				},
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Length": []string{"1000"},
				},
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.policy.IsCacheable(tt.response)
			if result != tt.expectedResult {
				t.Errorf("IsCacheable() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestPolicy_TTL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		policy         *pagecache.Policy
		response       *http.Response
		expectedResult time.Duration
	}{
		{
			name: "TTL with DefaultTTL",
			policy: func() *pagecache.Policy {
				p := pagecache.DefaultPolicy()
				p.UseCacheControl = false
				p.DefaultTTL = 30 * time.Minute
				return p
			}(),
			response: &http.Response{
				Header: http.Header{},
			},
			expectedResult: 30 * time.Minute,
		},
		{
			name: "TTL with Cache-Control max-age",
			policy: func() *pagecache.Policy {
				p := pagecache.DefaultPolicy()
				p.UseCacheControl = true
				return p
			}(),
			response: &http.Response{
				Header: http.Header{
					"Cache-Control": []string{"max-age=3600"},
				},
			},
			expectedResult: 3600 * time.Second,
		},
		{
			name: "TTL with Cache-Control max-age and invalid value",
			policy: func() *pagecache.Policy {
				p := pagecache.DefaultPolicy()
				p.UseCacheControl = true
				return p
			}(),
			response: &http.Response{
				Header: http.Header{
					"Cache-Control": []string{"max-age=invalid"},
				},
			},
			expectedResult: pagecache.DefaultTTL,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.policy.TTL(tt.response)
			if result != tt.expectedResult {
				t.Errorf("Expected TTL: %v, got: %v", tt.expectedResult, result)
			}
		})
	}
}

func parseTestURL(t *testing.T, urlStr string) *url.URL {
	t.Helper()

	uri, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}

	return uri
}
