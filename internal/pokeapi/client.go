package pokeapi

import (
	"io"
	"net/http"
	"time"

	"shiva.ddns.net/marcusvinicius/pokedex/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	pokecache  *pokecache.Cache
}

func NewClient(timeout, cacheInterval time.Duration) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		pokecache: pokecache.NewCache(cacheInterval),
	}
}

func (c *Client) makeGetRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
