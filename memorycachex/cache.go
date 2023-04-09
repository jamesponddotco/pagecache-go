package memorycachex

import (
	"context"
	"net/http"
	"sort"
	"sync"
	"time"

	"git.sr.ht/~jamesponddotco/pagecache-go"
)

// MemoryCache is an in-memory cache implementing the Cache interface.
type MemoryCache struct {
	cache       map[string]*Entry
	policy      *pagecache.Policy
	capacity    uint64
	currentSize uint64
	mu          sync.RWMutex
}

// Compile-time check to ensure Cache implements the cachex.Cache interface.
var _ pagecache.Cache = (*MemoryCache)(nil)

// NewCache creates a new MemoryCache instance with the specified policy and capacity.
func NewCache(policy *pagecache.Policy, capacity uint64) *MemoryCache {
	if policy == nil {
		policy = pagecache.DefaultPolicy()
	}

	if capacity <= 0 {
		capacity = pagecache.DefaultCapacity
	}

	return &MemoryCache{
		cache:    make(map[string]*Entry, capacity),
		policy:   policy,
		capacity: capacity,
		mu:       sync.RWMutex{},
	}
}

func (mc *MemoryCache) Get(_ context.Context, key string) (*http.Response, error) {
	mc.mu.RLock()
	entry, found := mc.cache[key]
	mc.mu.RUnlock()

	if !found {
		return nil, pagecache.ErrCacheMiss
	}

	if entry.IsExpired() {
		mc.mu.Lock()
		delete(mc.cache, key)
		mc.mu.Unlock()

		return nil, pagecache.ErrCacheMiss
	}

	entry.Access()

	response, err := entry.Load(key)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (mc *MemoryCache) Set(_ context.Context, key string, response *http.Response, expiration time.Duration) error {
	if !mc.policy.IsCacheable(response) { //nolint:contextcheck // we don't actually use the context for this package
		return nil
	}

	entry, err := NewEntry(key, response, time.Now().Add(expiration))
	if err != nil {
		return err
	}

	mc.mu.Lock()
	mc.cache[key] = entry
	mc.mu.Unlock()

	mc.evict()

	return nil
}

func (mc *MemoryCache) Delete(_ context.Context, key string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if _, found := mc.cache[key]; !found {
		return pagecache.ErrCacheMiss
	}

	delete(mc.cache, key)

	return nil
}

func (mc *MemoryCache) Policy() *pagecache.Policy {
	return mc.policy
}

func (mc *MemoryCache) Purge(_ context.Context) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.cache = make(map[string]*Entry)

	return nil
}

func (mc *MemoryCache) evict() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if mc.capacity == 0 {
		return
	}

	// Check if the cache is full or any entry is expired.
	var expired []*Entry

	for _, entry := range mc.cache {
		if entry.IsExpired() {
			expired = append(expired, entry)
			continue
		}

		if mc.currentSize < mc.capacity {
			return
		}

		mc.currentSize -= entry.Size
		delete(mc.cache, entry.Key)
	}

	// Evict entries based on the Mockingjay cache replacement policy.
	if len(expired) > 0 {
		sort.Slice(expired, func(i, j int) bool {
			return expired[i].Frequency < expired[j].Frequency
		})

		for _, entry := range expired {
			mc.currentSize -= entry.Size
			delete(mc.cache, entry.Key)
		}
	}
}
