package main

import (
	"fmt"
	"time"

	pokeapi "shiva.ddns.net/marcusvinicius/pokedex/internal"
)

func main() {
	pokeClient := pokeapi.NewClient(5 * time.Second)
	config := config{
		commands:      getCommands(),
		pokeapiClient: pokeClient,
	}

	if err := startRepl(&config); err != nil {
		fmt.Println(err)
	}
}
