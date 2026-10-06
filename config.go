package main

import (
	"shiva.ddns.net/marcusvinicius/pokedex/internal/pokeapi"
)

type config struct {
	commands         map[string]cliCommand
	caughtPokemons   map[string]pokeapi.PokemonData
	pokeapiClient    *pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
}
