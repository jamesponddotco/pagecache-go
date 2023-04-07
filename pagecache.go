// Package cachex implements a cache of HTTP responses.
package pagecache

import (
	"net/http"
	"strings"

	"git.sr.ht/~jamesponddotco/recache-go"
	"git.sr.ht/~jamesponddotco/recache-go/lrure"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"git.sr.ht/~jamesponddotco/xstd-go/xhash/xfnv"
)

const (
	// DefaultCacheName is the default name used for the default cache.
	DefaultCacheName string = "httpx"

	// DefaultCapacity is the default capacity of the memory cache when not
	// specified.
	DefaultCapacity uint64 = 128
)

// _RegexCache is a cache for regular expressions.
var _regexCache = lrure.New(recache.DefaultCapacity) //nolint:gochecknoglobals // global cache

// Common error messages for Cache implementations.
var (
	// ErrCacheMiss is returned when a cache entry is not found for the given key.
	ErrCacheMiss = NewCacheError(ErrNotFound, xerrors.Error("cache miss"))

	// ErrCacheExpired is returned when a cache entry is found but has expired.
	ErrCacheExpired = NewCacheError(ErrNotFound, xerrors.Error("cache entry expired"))

	// ErrKeyNotFound is returned when a cache entry is not found for the given key.
	ErrKeyNotFound = NewCacheError(ErrNotFound, xerrors.Error("key not found"))

	// ErrCacheStoreFailed is returned when storing an item in the cache fails.
	ErrCacheStoreFailed = NewCacheError(ErrOperationFailed, xerrors.Error("failed to set cache entry"))

	// ErrCacheDeleteFailed is returned when deleting an item from the cache fails.
	ErrCacheDeleteFailed = NewCacheError(ErrOperationFailed, xerrors.Error("failed to delete cache entry"))

	// ErrCachePurgeFailed is returned when purging the entire cache fails.
	ErrCachePurgeFailed = NewCacheError(ErrOperationFailed, xerrors.Error("failed to purge cache"))
)

// Key generates a cache key by concatenating information from the given
// *http.Request, cache name, and optional extra information. It hashes the
// result with the FNV-1 64-bit algorithm for fast hashing.
//
// The generated key is of the form:
//
//	"name:NAME:method:METHOD:url:URL:extra:EXTRA:EXTRA".
//
// This function is not used by the package itself, but is exported for use by
// packages implementing the Cache interface.
func Key(name string, req *http.Request, extra ...string) string {
	if name == strings.TrimSpace("") {
		name = DefaultCacheName
	}

	var builder strings.Builder

	builder.Grow(len(name) + len(req.Method) + len(req.URL.Redacted()) + len(extra)*2)

	builder.WriteString("name:")
	builder.WriteString(name)
	builder.WriteString(":method:")
	builder.WriteString(req.Method)
	builder.WriteString(":url:")
	builder.WriteString(req.URL.Redacted())

	if len(extra) > 0 {
		builder.WriteString(":extra:")

		for _, e := range extra {
			builder.WriteByte(':')
			builder.WriteString(e)
		}
	}

	return xfnv.String(builder.String())
}
