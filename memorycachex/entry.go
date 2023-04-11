package memorycachex

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"git.sr.ht/~jamesponddotco/pagecache-go"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

const (
	ErrKeyEmpty          xerrors.Error = "key must not be empty"
	ErrKeyMismatch       xerrors.Error = "key mismatch"
	ErrValueEmpty        xerrors.Error = "value must not be empty"
	ErrExpirationZero    xerrors.Error = "expiration must not be zero"
	ErrMarshalResponse   xerrors.Error = "failed to marshal response"
	ErrUnmarshalResponse xerrors.Error = "failed to unmarshal response"
)

// Entry represents a single cache entry. Entry is not tread-safe and should be
// protected by a sync.Mutex.
type Entry struct {
	Key        string
	Expiration time.Time
	Request    []byte
	Response   []byte
	Size       uint64
	Frequency  uint64
}

// Compile-time check to ensure Entry implements the cachex.Entry interface.
var _ pagecache.Entry = (*Entry)(nil)

// NewEntry creates a new cache entry with the specified key and expiration.
// You should call SerializeResponse after creating the entry to store the
// response.
func NewEntry(key string, resp *http.Response, expiration time.Time) (*Entry, error) {
	if key == "" {
		return nil, ErrKeyEmpty
	}

	if resp == nil {
		return nil, ErrValueEmpty
	}

	if expiration.IsZero() {
		return nil, ErrExpirationZero
	}

	request, response, err := pagecache.SaveResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrMarshalResponse, err)
	}

	entry := &Entry{
		Key:        key,
		Expiration: expiration,
		Request:    request,
		Response:   response,
		Size:       0,
		Frequency:  0,
	}

	return entry, nil
}

// Load loads the HTTP response from the cache entry.
func (e *Entry) Load(key string) (*http.Response, error) {
	if key != e.Key {
		return nil, ErrKeyMismatch
	}

	resp, err := pagecache.LoadResponse(e.Request, e.Response)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnmarshalResponse, err)
	}

	return resp, nil
}

// Access increments the frequency counter when the entry is accessed.
func (e *Entry) Access() {
	atomic.AddUint64(&e.Frequency, 1)
}

// SetSize updates the size of the cache entry.
func (e *Entry) SetSize(size uint64) {
	atomic.StoreUint64(&e.Size, size)
}

// SetTTL updates the expiration time of the cache entry.
func (e *Entry) SetTTL(ttl time.Duration) {
	e.Expiration = time.Now().Add(ttl)
}

// Expired checks if the cache entry has expired.
func (e *Entry) IsExpired() bool {
	if e.Expiration.IsZero() {
		return false
	}

	return time.Now().After(e.Expiration)
}
