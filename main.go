package main

import (
	"fmt"
	"time"

	"shiva.ddns.net/marcusvinicius/pokedex/internal/pokeapi"
)

const (
	requestTimeout = 5 * time.Second
	cacheInterval  = 10 * time.Minute
)

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
