package pokeapi

import (
	"encoding/json"
	"fmt"
)

// ListLocationEncounters -
func (c *Client) ListEncounters(locationName *string) (RespShallowEncounters, error) {
	url := baseURL + "/location-area/" + *locationName
	if locationName == nil {
		return RespShallowEncounters{}, fmt.Errorf("location name can't be blank")
	}

	if cachedEncounters, _ := c.pokecache.Get(url); cachedEncounters != nil {
		encountersResp := RespShallowEncounters{}
		_ = json.Unmarshal(cachedEncounters, &encountersResp)
		return encountersResp, nil
	}

	data, err := c.makeGetRequest(url)
	if err != nil {
		return RespShallowEncounters{}, err
	}

	encountersResp := RespShallowEncounters{}
	err = json.Unmarshal(data, &encountersResp)
	if err != nil {
		return RespShallowEncounters{}, err
	}

	c.pokecache.Add(url, data)

	return encountersResp, nil
}
