package main

import (
	"fmt"
	"time"

	"shiva.ddns.net/marcusvinicius/pokedex/internal/pokeapi"
)

const (
	interval = 5 * time.Second
)

func main() {
	pokeClient := pokeapi.NewClient(interval)
	config := config{
		commands:      getCommands(),
		pokeapiClient: pokeClient,
	}

	if err := startRepl(&config); err != nil {
		fmt.Println(err)
	}
}
