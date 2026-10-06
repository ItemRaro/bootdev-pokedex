package pokeapi

import (
	"encoding/json"
	"fmt"
)

// GetPokemonData returns a pokemon's data merged with its species data.
// The merged result is cached under the pokemon URL.
func (c *Client) GetPokemonData(pokemonName *string) (PokemonData, error) {
	if pokemonName == nil {
		return PokemonData{}, fmt.Errorf("pokemon name can't be blank")
	}

	url := c.baseUrl + "/pokemon/" + *pokemonName

	if cachedPokemonData, _ := c.pokecache.Get(url); cachedPokemonData != nil {
		pokemonDataResp := PokemonData{}
		_ = json.Unmarshal(cachedPokemonData, &pokemonDataResp)
		return pokemonDataResp, nil
	}

	pokemonDataResp := PokemonData{}

	data, err := c.makeGetRequest(url)
	if err != nil {
		return PokemonData{}, err
	}

	if err := json.Unmarshal(data, &pokemonDataResp); err != nil {
		return PokemonData{}, err
	}

	// Species data (capture rate, legendary and mythical flags) lives in another endpoint.
	data, err = c.makeGetRequest(pokemonDataResp.Species.URL)
	if err != nil {
		return PokemonData{}, err
	}

	if err := json.Unmarshal(data, &pokemonDataResp); err != nil {
		return PokemonData{}, err
	}

	// Cache the merged struct, not the raw response.
	data, err = json.Marshal(pokemonDataResp)
	if err != nil {
		return PokemonData{}, err
	}

	c.pokecache.Add(url, data)

	return pokemonDataResp, nil
}
