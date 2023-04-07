package pagecache

import (
	"context"
	"net/http"
	"time"
)

// Cache is a storage mechanism used to store and retrieve HTTP responses for
// improved reliability and performance.
//
// Implementations can use various caching strategies such as in-memory,
// file-based, or distributed caches like Redis.
type Cache interface {
	// Get retrieves an *http.Response from the cache associated with the given
	// key.
	Get(ctx context.Context, key string) (*http.Response, error)

	// Set stores an *http.Response in the cache, associated with the given
	// key, and sets the expiration duration.
	Set(ctx context.Context, key string, resp *http.Response, duration time.Duration) error

	// Delete removes the cache entry associated with the given key.
	Delete(ctx context.Context, key string) error

	// Policy returns the cache policy.
	Policy() *Policy

	// Purge clears the entire cache.
	Purge(ctx context.Context) error
}
