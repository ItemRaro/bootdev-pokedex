package main

import (
	"shiva.ddns.net/marcusvinicius/pokedex/internal/pokeapi"
)

// config holds the state shared between REPL commands.
type config struct {
	commands         map[string]cliCommand
	caughtPokemons   map[string]pokeapi.PokemonData
	pokeapiClient    *pokeapi.Client
	nextLocationsURL *string // next page for "map"; nil on the last page
	prevLocationsURL *string // previous page for "mapb"; nil on the first page
}
