package memorycachex_test

import (
	"reflect"
	"testing"

	"git.sr.ht/~jamesponddotco/pagecache-go"
	"git.sr.ht/~jamesponddotco/pagecache-go/memorycachex"
)

func TestNewCache(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give struct {
			policy   *pagecache.Policy
			capacity uint64
		}
		want struct {
			policy   *pagecache.Policy
			capacity uint64
		}
	}{
		{
			name: "Default policy and capacity",
			give: struct {
				policy   *pagecache.Policy
				capacity uint64
			}{
				policy:   nil,
				capacity: 0,
			},
			want: struct {
				policy   *pagecache.Policy
				capacity uint64
			}{
				policy:   pagecache.DefaultPolicy(),
				capacity: pagecache.DefaultCapacity,
			},
		},
		{
			name: "Custom policy and capacity",
			give: struct {
				policy   *pagecache.Policy
				capacity uint64
			}{
				policy:   pagecache.DefaultPolicy(),
				capacity: 100,
			},
			want: struct {
				policy   *pagecache.Policy
				capacity uint64
			}{
				policy:   pagecache.DefaultPolicy(),
				capacity: 100,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache := memorycachex.NewCache(tt.give.policy, tt.give.capacity)

			if !reflect.DeepEqual(cache.Policy(), tt.want.policy) {
				t.Errorf("Expected policy to be %v, but got %v", tt.want.policy, cache.Policy())
			}

			// Using reflection to access the unexported field "capacity"
			capacityField := reflect.ValueOf(cache).Elem().FieldByName("capacity")
			capacity := capacityField.Uint()

			if capacity != tt.want.capacity {
				t.Errorf("Expected capacity to be %v, but got %v", tt.want.capacity, capacity)
			}
		})
	}
}
