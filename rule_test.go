package pagecache_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/pagecache-go"
)

func TestRule_Match(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		rule           *pagecache.Rule
		url            string
		expectedResult bool
	}{
		{
			name: "Match with URL exact match",
			rule: &pagecache.Rule{
				URL:      "https://example.com/test",
				Behavior: pagecache.BehaviorInclude,
			},
			url:            "https://example.com/test",
			expectedResult: true,
		},
		{
			name: "Match with URL no match",
			rule: &pagecache.Rule{
				URL:      "https://example.com/test",
				Behavior: pagecache.BehaviorInclude,
			},
			url:            "https://example.com/other",
			expectedResult: false,
		},
		{
			name: "Match with regex pattern",
			rule: &pagecache.Rule{
				Pattern:  `^https://example\.com/\w+/\d+$`,
				Behavior: pagecache.BehaviorInclude,
			},
			url:            "https://example.com/words/123",
			expectedResult: true,
		},
		{
			name: "Match with regex pattern no match",
			rule: &pagecache.Rule{
				Pattern:  `^https://example\.com/\w+/\d+$`,
				Behavior: pagecache.BehaviorInclude,
			},
			url:            "https://example.com/words/abc",
			expectedResult: false,
		},
		{
			name: "Match with no URL and no pattern",
			rule: &pagecache.Rule{
				Behavior: pagecache.BehaviorInclude,
			},
			url:            "https://example.com/anything",
			expectedResult: false,
		},
		{
			name: "Match with invalid pattern",
			rule: &pagecache.Rule{
				Pattern:  "[invalid",
				Behavior: pagecache.BehaviorInclude,
			},
			url:            "https://example.com/anything",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.rule.Match(tt.url)
			if result != tt.expectedResult {
				t.Errorf("Expected Match: %v, got: %v", tt.expectedResult, result)
			}
		})
	}
}
