package main

import pokeapi "shiva.ddns.net/marcusvinicius/pokedex/internal"

type config struct {
	commands         map[string]cliCommand
	pokeapiClient    pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
}
