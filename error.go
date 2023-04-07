package pagecache

import "fmt"

const (
	// ErrNotFound indicates that the requested key was not found in the cache.
	ErrNotFound ErrorType = iota

	// ErrOperationFailed indicates that a cache operation, such as set or delete, has failed.
	ErrOperationFailed
)

// ErrorType represents the type of cache error.
type ErrorType int

// String returns a string representation of the ErrorType.
func (et ErrorType) String() string {
	switch et {
	case ErrNotFound:
		return "not found"
	case ErrOperationFailed:
		return "operation failed"
	default:
		return "unknown"
	}
}

// Error is a custom error type for the cache package.
type Error struct {
	// err is the underlying error.
	err error

	// errType is the type of error.
	errType ErrorType
}

// NewCacheError creates a new Error with the specified error type and underlying error.
func NewCacheError(errType ErrorType, err error) *Error {
	return &Error{
		err:     err,
		errType: errType,
	}
}

// Error returns a string representation of the CacheError.
func (ce *Error) Error() string {
	return fmt.Sprintf("cache error (%s): %v", ce.errType.String(), ce.err)
}

// Type returns the CacheErrorType of the CacheError.
func (ce *Error) Type() ErrorType {
	return ce.errType
}

// Unwrap returns the underlying error of the CacheError.
func (ce *Error) Unwrap() error {
	return ce.err
}
