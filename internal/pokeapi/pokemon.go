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

	data, err := c.makeGetRequest(url)
	if err != nil {
		return PokemonData{}, nil
	}

	shallowPokemonDataResp := RespShallowPokemonData{}
	if err := json.Unmarshal(data, &shallowPokemonDataResp); err != nil {
		return PokemonData{}, err
	}

	data, err = c.makeGetRequest(shallowPokemonDataResp.Species.URL)
	if err != nil {
		return PokemonData{}, nil
	}
	shallowPokemonSpeciesResp := RespShallowSpeciesData{}
	if err := json.Unmarshal(data, &shallowPokemonSpeciesResp); err != nil {
		return PokemonData{}, err
	}

	pokemonDataResp := PokemonData{
		shallowPokemonDataResp,
		shallowPokemonSpeciesResp,
	}

	data, err = json.Marshal(pokemonDataResp)
	if err != nil {
		return PokemonData{}, nil
	}

	c.pokecache.Add(url, data)

	return pokemonDataResp, nil
}
