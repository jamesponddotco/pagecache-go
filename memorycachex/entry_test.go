package memorycachex_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"git.sr.ht/~jamesponddotco/pagecache-go/memorycachex"
)

func TestNewEntry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		key         string
		resp        *http.Response
		expiration  time.Time
		expectError bool
	}{
		{
			name:        "Valid Entry",
			key:         "testkey",
			resp:        createValidResponse(t),
			expiration:  time.Now().Add(10 * time.Minute),
			expectError: false,
		},
		{
			name:        "Empty Key",
			key:         "",
			resp:        createValidResponse(t),
			expiration:  time.Now().Add(10 * time.Minute),
			expectError: true,
		},
		{
			name:        "Nil Response",
			key:         "testkey",
			resp:        nil,
			expiration:  time.Now().Add(10 * time.Minute),
			expectError: true,
		},
		{
			name:        "Zero Expiration",
			key:         "testkey",
			resp:        createValidResponse(t),
			expiration:  time.Time{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			entry, err := memorycachex.NewEntry(tt.key, tt.resp, tt.expiration)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if entry == nil {
					t.Errorf("Expected a non-nil entry, but got nil")
				} else {
					if entry.Key != tt.key {
						t.Errorf("Expected key %q, but got %q", tt.key, entry.Key)
					}
					if !entry.Expiration.Equal(tt.expiration) {
						t.Errorf("Expected expiration %v, but got %v", tt.expiration, entry.Expiration)
					}
				}
			}
		})
	}
}

func TestEntry_Load(t *testing.T) {
	t.Parallel()

	entry, err := createValidEntry(t)
	if err != nil {
		t.Fatalf("Failed to create a valid entry: %v", err)
	}

	tests := []struct {
		name        string
		key         string
		expectError bool
	}{
		{
			name:        "Valid Key",
			key:         "testkey",
			expectError: false,
		},
		{
			name:        "Invalid Key",
			key:         "wrongkey",
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp, err := entry.Load(tt.key)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if resp == nil {
					t.Errorf("Expected a non-nil response, but got nil")
				} else if resp.StatusCode != http.StatusOK {
					t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
				}
			}
		})
	}
}

func TestEntry_Access(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give uint64
		want uint64
	}{
		{
			name: "Access Once",
			give: 1,
			want: 1,
		},
		{
			name: "Access Multiple Times",
			give: 5,
			want: 5,
		},
		{
			name: "Access Zero Times",
			give: 0,
			want: 0,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			entry, err := createValidEntry(t)
			if err != nil {
				t.Fatalf("Failed to create a valid entry: %v", err)
			}

			for i := uint64(0); i < tt.give; i++ {
				entry.Access()
			}

			actualCount := entry.Frequency
			if actualCount != tt.want {
				t.Errorf("Expected access count to be %d, but got %d", tt.want, actualCount)
			}
		})
	}
}

func TestEntry_SetSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size uint64
	}{
		{
			name: "SetSize to 0",
			size: 0,
		},
		{
			name: "SetSize to a small value",
			size: 1024,
		},
		{
			name: "SetSize to a large value",
			size: 1024 * 1024 * 1024,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			entry, err := createValidEntry(t)
			if err != nil {
				t.Fatalf("Failed to create a new entry: %v", err)
			}

			entry.SetSize(tt.size)

			if entry.Size != tt.size {
				t.Errorf("Expected size to be %d, but got %d", tt.size, entry.Size)
			}
		})
	}
}

func TestEntry_SetTTL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ttl  time.Duration
	}{
		{
			name: "SetTTL to a small value",
			ttl:  5 * time.Second,
		},
		{
			name: "SetTTL to a medium value",
			ttl:  30 * time.Minute,
		},
		{
			name: "SetTTL to a large value",
			ttl:  48 * time.Hour,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			entry, err := createValidEntry(t)
			if err != nil {
				t.Fatalf("Failed to create a new entry: %v", err)
			}

			entry.SetTTL(tt.ttl)
			newExpiration := time.Now().Add(tt.ttl)
			if entry.Expiration.Sub(newExpiration) > time.Millisecond {
				t.Errorf("Expected expiration time to be near %v, but got %v", newExpiration, entry.Expiration)
			}
		})
	}
}

func TestEntry_IsExpired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		expiration time.Time
		want       bool
	}{
		{
			name:       "Not expired",
			expiration: time.Now().Add(5 * time.Minute),
			want:       false,
		},
		{
			name:       "Expired",
			expiration: time.Now().Add(-5 * time.Minute),
			want:       true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			entry, err := memorycachex.NewEntry("testkey", createValidResponse(t), tt.expiration)
			if err != nil {
				t.Fatalf("Failed to create a new entry: %v", err)
			}

			got := entry.IsExpired()
			if got != tt.want {
				t.Errorf("Expected IsExpired() to be %v, but got %v", tt.want, got)
			}
		})
	}
}

func createValidResponse(t *testing.T) *http.Response {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, "http://example.com/", http.NoBody)
	if err != nil {
		panic(err)
	}

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusOK)
	rec.WriteString("OK")

	resp := rec.Result()
	resp.Request = req

	return resp
}

func createValidEntry(t *testing.T) (*memorycachex.Entry, error) {
	t.Helper()

	var (
		resp       = createValidResponse(t)
		expiration = time.Now().Add(10 * time.Minute)
	)

	return memorycachex.NewEntry("testkey", resp, expiration)
}
