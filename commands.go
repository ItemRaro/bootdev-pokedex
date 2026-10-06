package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"

	"shiva.ddns.net/marcusvinicius/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func commandMapf(cfg *config, args ...string) error {
	locationsResp, err := cfg.pokeapiClient.ListLocations(cfg.nextLocationsURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locationsResp.Next
	cfg.prevLocationsURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.prevLocationsURL == nil {
		return errors.New("you're on the first page")
	}

	locationResp, err := cfg.pokeapiClient.ListLocations(cfg.prevLocationsURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locationResp.Next
	cfg.prevLocationsURL = locationResp.Previous

	for _, loc := range locationResp.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("location name can't be empty")
	}

	locationName := args[0]
	encountersResp, err := cfg.pokeapiClient.ListEncounters(&locationName)
	if err != nil {
		return err
	}

	for _, encounter := range encountersResp.Encounters {
		fmt.Println(encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("location name can't be empty")
	}

	pokemonName := args[0]
	pokemonData, err := cfg.pokeapiClient.GetPokemonData(&pokemonName)
	if err != nil {
		return err
	}

	fmt.Println("Throwing a Pokeball at " + pokemonData.Name + "...")
	if catchPokemon(pokemonData.CaptureRate, pokemonData.IsLegendary, pokemonData.IsMythical) {
		fmt.Println(pokemonData.Name + " was caught!")
		cfg.caughtPokemons[pokemonData.Name] = pokemonData

		return nil
	} else {
		fmt.Println(pokemonData.Name + " escaped!")

		return nil
	}
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("pokemon name can't be empty")
	}

	pokemonName := args[0]
	pokemon, ok := cfg.caughtPokemons[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	inspectPokemon(pokemon)

	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	if len(cfg.caughtPokemons) <= 0 {
		fmt.Println("your Pokedex is empty! Catch some Pokemons!")

		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, pokemon := range cfg.caughtPokemons {
		fmt.Printf("  -%s\n", pokemon.Name)
	}

	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range cfg.commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()

	return nil
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")

	os.Exit(0)

	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Get the next page of locations",
			callback:    commandMapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous page of locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location. Use explore + location-name",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon. Use catch + pokemon-name",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect caught pokemon. Use inspect + pokemon-name",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists your capture pokemons",
			callback:    commandPokedex,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func catchPokemon(captureRate int, legendary, mythical bool) bool {
	const (
		captureRateCap = 255
		legendaryCap   = 45
		mythicalCap    = 30
	)

	if legendary && captureRate > legendaryCap {
		captureRate = legendaryCap
	}

	if mythical && captureRate > mythicalCap {
		captureRate = mythicalCap
	}

	roll := rand.IntN(captureRateCap + 1)

	return roll < captureRate
}

func inspectPokemon(pokemonData pokeapi.PokemonData) {
	fmt.Println("Name: " + pokemonData.Name)
	fmt.Printf("Height: %d\n", pokemonData.Height)
	fmt.Printf("Weight: %d\n", pokemonData.Weight)
	fmt.Println("Stats:")
	for i := range pokemonData.Stats {
		fmt.Printf("  -%s: %d\n", pokemonData.Stats[i].Stat.Name, pokemonData.Stats[i].BaseStat)
	}
	fmt.Println("Types:")
	for i := range pokemonData.Types {
		fmt.Printf("  -%s\n", pokemonData.Types[i].Type.Name)
	}
}
