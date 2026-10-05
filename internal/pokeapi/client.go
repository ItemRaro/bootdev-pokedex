package pokeapi

import (
	"net/http"
	"time"

	"shiva.ddns.net/marcusvinicius/pokedex/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	pokecache  *pokecache.Cache
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		pokecache: pokecache.NewCache(timeout),
	}
}
