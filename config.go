package main

import (
	"shiva.ddns.net/marcusvinicius/pokedex/internal/pokeapi"
)

type config struct {
	commands         map[string]cliCommand
	pokeapiClient    *pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
}
