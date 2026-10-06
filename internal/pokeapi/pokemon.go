package pokeapi

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetPokemonData(pokemonName *string) (PokemonData, error) {
	url := baseURL + "/pokemon/" + *pokemonName

	if pokemonName == nil {
		return PokemonData{}, fmt.Errorf("pokemon name can't be blank")
	}

	if cachedPokemonData, _ := c.pokecache.Get(url); cachedPokemonData != nil {
		pokemonDataResp := PokemonData{}
		_ = json.Unmarshal(cachedPokemonData, &pokemonDataResp)
		return pokemonDataResp, nil
	}

	pokemonDataResp := PokemonData{}

	data, err := c.makeGetRequest(url)
	if err != nil {
		return PokemonData{}, nil
	}

	if err := json.Unmarshal(data, &pokemonDataResp); err != nil {
		return PokemonData{}, err
	}

	data, err = c.makeGetRequest(pokemonDataResp.Species.URL)
	if err != nil {
		return PokemonData{}, nil
	}

	if err := json.Unmarshal(data, &pokemonDataResp); err != nil {
		return PokemonData{}, err
	}

	data, err = json.Marshal(pokemonDataResp)
	if err != nil {
		return PokemonData{}, nil
	}

	c.pokecache.Add(url, data)

	return pokemonDataResp, nil
}
