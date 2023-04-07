package pagecache

import (
	"context"

	"git.sr.ht/~jamesponddotco/recache-go"
)

const (
	// BehaviorInclude means to include the URL in caching, even if it would be
	// excluded by default.
	BehaviorInclude Behavior = iota

	// BehaviorExclude means to exclude the URL from caching.
	BehaviorExclude
)

// Behavior represents the caching behavior for a specific URL pattern.
type Behavior int

// Rule defines a pattern for matching URLs, a caching behavior, and an
// optional map of custom headers to apply when the rule is matched.
type Rule struct {
	// URL is a URL to match against URLs.
	URL string

	// Pattern is a regular expression pattern to match against URLs.
	//
	// This field is ignored if URL is set.
	Pattern string

	// PatternFlag is a control flag for the regular expression pattern.
	//
	// URL is a slice of URLs to match against URLs.
	PatternFlag recache.Flag

	// Behavior is the caching behavior to apply for the matched URLs.
	Behavior Behavior
}

// Match returns true if the URL matches the rule.
func (r *Rule) Match(url string) bool {
	if r.URL != "" {
		return r.URL == url
	}

	if r.Pattern != "" {
		re, err := _regexCache.Get(context.Background(), r.Pattern, r.PatternFlag)
		if err != nil {
			return false
		}

		return re.MatchString(url)
	}

	return false
}
