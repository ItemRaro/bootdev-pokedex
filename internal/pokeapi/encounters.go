package pokeapi

import (
	"encoding/json"
	"fmt"
)

// ListEncounters returns the pokemon that can be found in a location area.
func (c *Client) ListEncounters(locationName *string) (RespShallowEncounters, error) {
	if locationName == nil {
		return RespShallowEncounters{}, fmt.Errorf("location name can't be blank")
	}

	url := c.baseUrl + "/location-area/" + *locationName

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
