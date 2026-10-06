package main

import (
	"fmt"
	"time"

	"shiva.ddns.net/marcusvinicius/pokedex/internal/pokeapi"
)

const (
	requestTimeout = 5 * time.Second  // limit for each HTTP request
	cacheInterval  = 10 * time.Minute // how long API responses stay cached
)

// main wires up the API client and starts the REPL.
func main() {
	pokeClient := pokeapi.NewClient(requestTimeout, cacheInterval)
	config := config{
		commands:       getCommands(),
		caughtPokemons: make(map[string]pokeapi.PokemonData),
		pokeapiClient:  pokeClient,
	}

	if err := startRepl(&config); err != nil {
		fmt.Println(err)
	}
}
