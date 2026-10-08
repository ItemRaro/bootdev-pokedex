// Package pokeapi is a small client for PokeAPI (https://pokeapi.co) with response caching.
package pokeapi

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"shiva.ddns.net/marcusvinicius/pokedex/internal/pokecache"
)

// Client performs HTTP requests to PokeAPI and caches the responses.
type Client struct {
	baseUrl    string
	httpClient http.Client
	pokecache  *pokecache.Cache
}

// NewClient creates a client with the given request timeout and cache expiration.
func NewClient(timeout, cacheInterval time.Duration) *Client {
	return &Client{
		baseUrl: "https://pokeapi.co/api/v2",
		httpClient: http.Client{
			Timeout: timeout,
		},
		pokecache: pokecache.NewCache(cacheInterval),
	}
}

// makeGetRequest performs a GET request and returns the raw response body.
func (c *Client) makeGetRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Anything other than 200 (e.g. 404 for an unknown pokemon) is an error.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
