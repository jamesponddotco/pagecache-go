package pagecache

import (
	"net/http"
	"time"
)

// Entry is an interface that represents a cache entry.
type Entry interface {
	// Load loads the HTTP response from the cache entry.
	Load(key string) (*http.Response, error)

	// Access updates the last access time of the entry.
	Access()

	// SetTTL sets the time-to-live of the entry.
	SetTTL(ttl time.Duration)

	// SetSize sets the size of the entry.
	SetSize(size uint64)

	// IsExpired checks if the entry is expired.
	IsExpired() bool
}
