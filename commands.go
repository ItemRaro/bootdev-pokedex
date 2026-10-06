package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"

	"shiva.ddns.net/marcusvinicius/pokedex/internal/pokeapi"
)

// cliCommand describes a REPL command and the function that runs it.
type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

// commandMapf prints the next page of location areas.
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

// commandMapb prints the previous page of location areas.
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

// commandExplore lists the pokemon found in a location area.
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

// commandCatch tries to catch a pokemon and adds it to the pokedex on success.
func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("pokemon name can't be empty")
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

// commandInspect shows the details of a caught pokemon.
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

// commandPokedex lists every caught pokemon.
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

// commandHelp prints the available commands.
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

// commandExit closes the program.
func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")

	os.Exit(0)

	return nil
}

// getCommands returns every command available in the REPL, keyed by name.
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

// catchPokemon rolls a 0-255 die against the capture rate (higher is easier).
// Legendary and mythical pokemon have capped rates so they are never too easy.
func catchPokemon(captureRate int, legendary, mythical bool) bool {
	const (
		captureRateCap = 255 // highest capture rate in PokeAPI
		legendaryCap   = 45  // ~17.6% max chance
		mythicalCap    = 30  // ~11.7% max chance
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

// inspectPokemon prints a pokemon's name, size, stats and types.
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
