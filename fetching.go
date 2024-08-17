package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type stash struct {
	createdAt time.Time
	val       []byte
}

var cacheMutex sync.RWMutex

func cleanCache(threshold time.Duration) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	now := time.Now()
	for url, s := range responseCache {
		if now.Sub(s.createdAt) > threshold {
			delete(responseCache, url)
		}
	}
}

// the use of RWMutex is not that of necessity, perhaps
func FetchUnmarshalCache[T any](url string) (T, error) {
	var result T

	cacheMutex.RLock()
	// fmt.Println("cache is now of length:", len(responseCache))
	// fmt.Println("search for cache with key:", url)
	cachedBody, ok := responseCache[url]
	cacheMutex.RUnlock()

	if ok {
		err := json.Unmarshal(cachedBody.val, &result)
		if err != nil {
			return result, fmt.Errorf("failed to unmarshal cached JSON: %w", err)
		}
		// fmt.Println("FROM CACHE---")
		return result, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return result, fmt.Errorf("failed to fetch data from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		cacheMutex.Lock()
		responseCache[url] = stash{time.Now(), body}
		cacheMutex.Unlock()
		return result, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	cacheMutex.Lock()
	// fmt.Println("STORING map of:", len(responseCache))
	// fmt.Println("STORING cache with key:", url)
	responseCache[url] = stash{time.Now(), body}

	// fmt.Printf("%+v \n", responseCache[url])

	cacheMutex.Unlock()

	return result, nil
}

func locMap(s *state) error {

	endpoint := s.mapNext
	if endpoint == "" {
		endpoint = "http://pokeapi.co/api/v2/location"
	}

	locRes, err := FetchUnmarshalCache[locationResponse](endpoint)
	if err != nil {
		return err
	}

	if !s.testing {
		for _, location := range locRes.Results {
			fmt.Println(location.Name)
		}
	}

	s.mapNext = locRes.Next
	s.mapBefore = endpoint

	return nil
}

func locMapB(s *state) error {

	endpoint := s.mapBefore
	if endpoint == "" {
		endpoint = "http://pokeapi.co/api/v2/location"
	}

	locRes, err := FetchUnmarshalCache[locationResponse](endpoint)
	if err != nil {
		return err
	}

	if !s.testing {
		for _, location := range locRes.Results {
			fmt.Println(location.Name)
		}
	}

	s.mapNext = locRes.Next
	s.mapBefore = locRes.Previous

	return nil
}

// relicts

// func FetchUnmarshalCache[T any](url string) (T, error) {
// 	var result T
// 	// var body []byte
// 	if cachedBody, ok := responseCache[url]; ok {
// 		// If found in cache, unmarshal the cached data
// 		err := json.Unmarshal(cachedBody.val, &result)
// 		if err != nil {
// 			return result, fmt.Errorf("failed to unmarshal cached JSON: %w", err)
// 		}
// 		return result, nil
// 	}

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return result, fmt.Errorf("failed to fetch data from URL: %w", err)
// 	}

// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return result, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return result, fmt.Errorf("failed to read response body: %w", err)
// 	}

// 	err = json.Unmarshal([]byte(body), &result)
// 	if err != nil {
// 		responseCache[url] = stash{time.Now(), []byte(body)}
// 		return result, fmt.Errorf("failed to unmarshal JSON: %w", err)
// 	}
// 	// fmt.Println("result should be returned", result)
// 	return result, nil
// }
