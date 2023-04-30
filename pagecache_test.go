package pagecache_test

import (
	"net/http"
	"testing"

	"git.sr.ht/~jamesponddotco/pagecache-go"
)

func TestKey(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodGet, "https://example.com/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		giveName  string
		giveExtra []string
		want      string
	}{
		{
			name:      "default cache name and some extras",
			giveName:  "",
			giveExtra: []string{"foo", "bar"},
			want:      "984a302b2754ae67",
		},
		{
			name:      "custom cache name and some extras",
			giveName:  "foo",
			giveExtra: []string{"foo", "bar"},
			want:      "e78d179014d38edf",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := pagecache.Key(tt.giveName, req, tt.giveExtra...)
			if got != tt.want {
				t.Errorf("Key() = %v, want %v", got, tt.want)
			}
		})
	}
}
