package sliceutil_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/pagecache-go/internal/sliceutil"
)

func TestMatchString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give struct {
			slice []string
			s     string
		}
		want bool
	}{
		{
			name: "match found",
			give: struct {
				slice []string
				s     string
			}{
				slice: []string{"apple", "banana", "cherry"},
				s:     "banana",
			},
			want: true,
		},
		{
			name: "match not found",
			give: struct {
				slice []string
				s     string
			}{
				slice: []string{"apple", "banana", "cherry"},
				s:     "grape",
			},
			want: false,
		},
		{
			name: "empty slice",
			give: struct {
				slice []string
				s     string
			}{
				slice: []string{},
				s:     "banana",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := sliceutil.MatchString(tt.give.slice, tt.give.s); got != tt.want {
				t.Errorf("MatchString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give struct {
			slice []int
			i     int
		}
		want bool
	}{
		{
			name: "match found",
			give: struct {
				slice []int
				i     int
			}{
				slice: []int{1, 2, 3},
				i:     2,
			},
			want: true,
		},
		{
			name: "match not found",
			give: struct {
				slice []int
				i     int
			}{
				slice: []int{1, 2, 3},
				i:     4,
			},
			want: false,
		},
		{
			name: "empty slice",
			give: struct {
				slice []int
				i     int
			}{
				slice: []int{},
				i:     2,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := sliceutil.MatchInt(tt.give.slice, tt.give.i); got != tt.want {
				t.Errorf("MatchInt() = %v, want %v", got, tt.want)
			}
		})
	}
}