package httputil_test

import (
	"net/http"
	"testing"

	"git.sr.ht/~jamesponddotco/pagecache-go/internal/httputil"
)

func TestIsBodySizeWithinLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		header http.Header
		limit  int64
		want   bool
	}{
		{
			name:   "Negative limit",
			header: http.Header{"Content-Length": []string{"100"}},
			limit:  -1,
			want:   true,
		},
		{
			name:   "Zero limit",
			header: http.Header{"Content-Length": []string{"100"}},
			limit:  0,
			want:   true,
		},
		{
			name:   "Content-Length within limit",
			header: http.Header{"Content-Length": []string{"100"}},
			limit:  200,
			want:   true,
		},
		{
			name:   "Content-Length exceeding limit",
			header: http.Header{"Content-Length": []string{"300"}},
			limit:  200,
			want:   false,
		},
		{
			name:   "Invalid Content-Length",
			header: http.Header{"Content-Length": []string{"invalid"}},
			limit:  200,
			want:   true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := httputil.IsBodySizeWithinLimit(tt.header, tt.limit)
			if got != tt.want {
				t.Errorf("IsBodySizeWithinLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxAge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		header http.Header
		want   int
	}{
		{
			name:   "Nil header",
			header: nil,
			want:   -1,
		},
		{
			name:   "No Cache-Control header",
			header: http.Header{"Content-Type": []string{"text/plain"}},
			want:   -1,
		},
		{
			name:   "Cache-Control with max-age",
			header: http.Header{"Cache-Control": []string{"private, max-age=300"}},
			want:   300,
		},
		{
			name:   "Cache-Control without max-age",
			header: http.Header{"Cache-Control": []string{"private, no-store"}},
			want:   -1,
		},
		{
			name:   "Cache-Control with invalid max-age",
			header: http.Header{"Cache-Control": []string{"private, max-age=invalid"}},
			want:   -1,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := httputil.MaxAge(tt.header)
			if got != tt.want {
				t.Errorf("MaxAge() = %v, want %v", got, tt.want)
			}
		})
	}
}
