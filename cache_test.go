package main

import (
	"testing"
)

// Test the FetchUnmarshalCache function
func TestFetchUnmarshalCache(t *testing.T) {
	responseCache = make(map[string]stash)
	cases := []string{"hello"}
	for _, url := range cases {
		FetchUnmarshalCache[locationResponse](url)
	}
}

// Test the cleanCache function
func TestCleanCache(t *testing.T) {
	responseCache = make(map[string]stash)

}
