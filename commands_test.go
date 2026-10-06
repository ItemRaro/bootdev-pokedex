package main

import (
	"testing"

	"shiva.ddns.net/marcusvinicius/pokedex/internal/pokeapi"
)

func TestCatchPokemonNeverCatchesZeroRate(t *testing.T) {
	for range 1000 {
		if catchPokemon(0, false, false) {
			t.Errorf("expected capture rate 0 to never catch")
			return
		}
	}
}

// TestCatchPokemonRates rolls many times and checks the success rate is close to
// capped rate / 256. The tolerance covers normal randomness.
func TestCatchPokemonRates(t *testing.T) {
	const (
		tries     = 100_000
		tolerance = 2.0 // percentage points
	)

	cases := []struct {
		name        string
		captureRate int
		legendary   bool
		mythical    bool
		expected    float64 // percent
	}{
		{
			name:        "common keeps its rate",
			captureRate: 128,
			expected:    128.0 / 256 * 100,
		},
		{
			name:        "legendary is capped",
			captureRate: 255,
			legendary:   true,
			expected:    45.0 / 256 * 100,
		},
		{
			name:        "mythical is capped",
			captureRate: 45,
			mythical:    true,
			expected:    30.0 / 256 * 100,
		},
		{
			name:        "legendary below cap keeps its rate",
			captureRate: 3,
			legendary:   true,
			expected:    3.0 / 256 * 100,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			caught := 0
			for range tries {
				if catchPokemon(c.captureRate, c.legendary, c.mythical) {
					caught++
				}
			}

			rate := float64(caught) / tries * 100
			if rate < c.expected-tolerance || rate > c.expected+tolerance {
				t.Errorf("expected ~%.1f%%, got %.1f%%", c.expected, rate)
				return
			}
		})
	}
}

// TestCommandErrors only covers paths that return before any API request.
func TestCommandErrors(t *testing.T) {
	cases := []struct {
		name      string
		command   func(*config, ...string) error
		args      []string
		wantError bool
	}{
		{
			name:      "mapb on the first page",
			command:   commandMapb,
			wantError: true,
		},
		{
			name:      "explore without a location",
			command:   commandExplore,
			wantError: true,
		},
		{
			name:      "catch without a pokemon",
			command:   commandCatch,
			wantError: true,
		},
		{
			name:      "inspect without a pokemon",
			command:   commandInspect,
			wantError: true,
		},
		{
			name:      "inspect a pokemon not caught",
			command:   commandInspect,
			args:      []string{"pidgey"},
			wantError: false,
		},
		{
			name:      "pokedex when empty",
			command:   commandPokedex,
			wantError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := &config{
				commands:       getCommands(),
				caughtPokemons: make(map[string]pokeapi.PokemonData),
			}

			err := c.command(cfg, c.args...)
			if (err != nil) != c.wantError {
				t.Errorf("expected error: %v, got: %v", c.wantError, err)
				return
			}
		})
	}
}
