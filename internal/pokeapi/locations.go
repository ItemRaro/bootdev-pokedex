package pokeapi

import (
	"encoding/json"
)

// ListLocations returns a page of location areas.
// A nil pageURL fetches the first page.
func (c *Client) ListLocations(pageURL *string) (RespShallowLocations, error) {
	url := c.baseUrl + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	if cachedLocations, _ := c.pokecache.Get(url); cachedLocations != nil {
		locationsResp := RespShallowLocations{}
		_ = json.Unmarshal(cachedLocations, &locationsResp)
		return locationsResp, nil
	}

	data, err := c.makeGetRequest(url)
	if err != nil {
		return RespShallowLocations{}, err
	}

	locationsResp := RespShallowLocations{}
	err = json.Unmarshal(data, &locationsResp)
	if err != nil {
		return RespShallowLocations{}, err
	}

	c.pokecache.Add(url, data)

	return locationsResp, nil
}
