package main

import (
	"fmt"
	"testing"
	"time"
)

// func TestUrlBuilder(t *testing.T) {
// 	baseURL := &url.URL{
// 		Scheme: "https",
// 		Host:   "example.com",
// 		Path:   "search",
// 	}

// 	// Set query parameters
// 	query := url.Values{}
// 	query.Set("q", "golang")
// 	query.Set("page", "1")

// 	// Encode the query parameters into the URL
// 	baseURL.RawQuery = query.Encode()

// 	fmt.Println(baseURL.String())
// }

// Test the FetchUnmarshalCache function
func TestFetchUnmarshalCache(t *testing.T) {
	responseCache = make(map[string]stash)
	defer func() { responseCache = make(map[string]stash) }()

	cases := []string{"http://pokeapi.co/api/v2/location"}
	fmt.Printf("requesting %v cases \n", len(cases))
	for _, urlx := range cases {
		FetchUnmarshalCache[locationResponse](urlx)
		if _, ok := responseCache[urlx]; !ok {
			t.Fatalf("the %s key is not there ", urlx)
		}
	}

}

func TestMapMapB(t *testing.T) {
	responseCache = make(map[string]stash)
	ss = state{0, 20, "", "", false, true}

	responseCache = make(map[string]stash)
	defer func() {
		responseCache = make(map[string]stash)
		ss = state{0, 20, "", "", false, true}
	}()

	cases := make(map[int]int)
	cases[3] = 3
	cases[0] = 0
	cases[5] = 5

	for key, val := range cases {
		if len(responseCache) != 0 {
			t.Fatalf("map should be 0 at start but its %v", len(responseCache))
		}
		for i := 0; i < key; i++ {
			locMap(&ss)
		}
		if len(responseCache) != val {
			t.Fatalf("map %v times, cache size should be %v but its %v", key, val, len(responseCache))
		}
		for i := 0; i < key; i++ {
			locMapB(&ss)
		}
		if len(responseCache) < val {
			t.Fatalf("mapBack %v times, cache size should greater than 	 %v but its %v", key, val, len(responseCache))
		}
		responseCache = make(map[string]stash)
		ss = state{0, 20, "", "", false, true}
	}
	// locMap(&ss)

}

// Test the cleanCache function
func TestCleanCache(t *testing.T) {
	responseCache = make(map[string]stash)
	ss = state{0, 20, "", "", false, true}
	// time.Sleep(time.Second * 10)
	// fmt.Println("slept 10 seconfs")

	cases := make(map[int]int)
	cases[3] = 1
	cases[0] = 0
	cases[5] = 1

	for key, val := range cases {
		if len(responseCache) != 0 {
			t.Fatalf("map should be 0 at start but its %v", len(responseCache))
		}
		for i := 0; i < key; i++ {
			if i == key-1 {
				time.Sleep(time.Second * 8)
			}
			locMap(&ss)
		}
		cleanCache(time.Second * 7)
		if len(responseCache) != val {
			t.Fatalf("req %v times then sleep, cache size should be %v but its %v", key, val, len(responseCache))
		}

		responseCache = make(map[string]stash)
		ss = state{0, 20, "", "", false, true}
	}

}
