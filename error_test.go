package pagecache_test

import (
	"fmt"
	"testing"

	"git.sr.ht/~jamesponddotco/pagecache-go"
)

func TestErrorType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		errType  pagecache.ErrorType
		expected string
	}{
		{
			name:     "ErrNotFound",
			errType:  pagecache.ErrNotFound,
			expected: "not found",
		},
		{
			name:     "ErrOperationFailed",
			errType:  pagecache.ErrOperationFailed,
			expected: "operation failed",
		},
		{
			name:     "UnknownErrorType",
			errType:  pagecache.ErrorType(999),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := tt.errType.String()
			if actual != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, actual)
			}
		})
	}
}

func TestCacheError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		errType   pagecache.ErrorType
		err       error
		expected  string
		unwrapped error
	}{
		{
			name:      "ErrNotFound",
			errType:   pagecache.ErrNotFound,
			err:       fmt.Errorf("key not found"),
			expected:  "cache error (not found): key not found",
			unwrapped: fmt.Errorf("key not found"),
		},
		{
			name:      "ErrOperationFailed",
			errType:   pagecache.ErrOperationFailed,
			err:       fmt.Errorf("set operation failed"),
			expected:  "cache error (operation failed): set operation failed",
			unwrapped: fmt.Errorf("set operation failed"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cacheErr := pagecache.NewCacheError(tt.errType, tt.err)

			if cacheErr.Error() != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, cacheErr.Error())
			}

			if cacheErr.Type() != tt.errType {
				t.Errorf("expected ErrorType '%d', got '%d'", tt.errType, cacheErr.Type())
			}

			if cacheErr.Unwrap().Error() != tt.unwrapped.Error() {
				t.Errorf("Unwrap() does not return the expected underlying error")
			}
		})
	}
}
