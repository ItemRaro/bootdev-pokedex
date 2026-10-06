package pokeapi

import (
	"encoding/json"
)

// ListLocations -
func (c *Client) ListLocations(pageURL *string) (RespShallowLocations, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		if cachedLocations, _ := c.pokecache.Get(*pageURL); cachedLocations != nil {
			locationsResp := RespShallowLocations{}
			_ = json.Unmarshal(cachedLocations, &locationsResp)
			return locationsResp, nil
		}
		url = *pageURL
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
