package pokeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// realAPI is the address written inside the fixtures' links (next, species.url...).
const realAPI = "https://pokeapi.co/api/v2"

// fixtures maps each request (path + query) to the testdata file that answers it.
// The species route uses the ID because it comes from species.url in pokemon_pidgey.json.
var fixtures = map[string]string{
	"/location-area":                    "location_area_page1.json",
	"/location-area?offset=20&limit=20": "location_area_page2.json",
	"/location-area/canalave-city-area": "location_area_canalave.json",
	"/pokemon/pidgey":                   "pokemon_pidgey.json",
	"/pokemon-species/16/":              "species_pidgey.json",
}

// newTestClient starts a fake PokeAPI serving the fixtures and returns a client
// pointed at it, plus the number of requests received per route.
func newTestClient(t *testing.T) (*Client, map[string]int) {
	t.Helper()

	calls := map[string]int{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := r.URL.RequestURI()
		calls[route]++

		if route == "/pokemon/broken" {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"name":"broken"}`))
			return
		}

		file, ok := fixtures[route]
		if !ok {
			http.NotFound(w, r)
			return
		}

		data, err := os.ReadFile("testdata/" + file)
		if err != nil {
			t.Errorf("missing fixture %s: %v", file, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Point the fixtures' links at this fake server instead of the real API.
		body := strings.ReplaceAll(string(data), realAPI, "http://"+r.Host)
		_, _ = w.Write([]byte(body))
	})

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	client := NewClient(5*time.Second, time.Minute)
	client.baseUrl = srv.URL

	return client, calls
}

func TestListLocations(t *testing.T) {
	client, _ := newTestClient(t)

	firstPage, err := client.ListLocations(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if len(firstPage.Results) != 20 {
		t.Errorf("expected 20 locations, got %d", len(firstPage.Results))
		return
	}
	if firstPage.Results[0].Name != "canalave-city-area" {
		t.Errorf("expected canalave-city-area, got %s", firstPage.Results[0].Name)
		return
	}
	if firstPage.Previous != nil {
		t.Errorf("expected no previous page on the first page")
		return
	}
	if firstPage.Next == nil {
		t.Errorf("expected a next page URL")
		return
	}

	// Follow the "next" link, like the map command does.
	secondPage, err := client.ListLocations(firstPage.Next)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if secondPage.Results[0].Name != "mt-coronet-1f-route-216" {
		t.Errorf("expected mt-coronet-1f-route-216, got %s", secondPage.Results[0].Name)
		return
	}
	if secondPage.Previous == nil {
		t.Errorf("expected a previous page URL")
		return
	}
}

func TestListEncounters(t *testing.T) {
	cases := []struct {
		location string
		expected []string
	}{
		{
			location: "canalave-city-area",
			expected: []string{
				"tentacool", "tentacruel", "staryu", "magikarp", "gyarados", "wingull",
				"pelipper", "shellos", "gastrodon", "finneon", "lumineon",
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			client, _ := newTestClient(t)

			resp, err := client.ListEncounters(&c.location)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if len(resp.Encounters) != len(c.expected) {
				t.Errorf("expected %d pokemon, got %d", len(c.expected), len(resp.Encounters))
				return
			}
			for j, encounter := range resp.Encounters {
				if encounter.Pokemon.Name != c.expected[j] {
					t.Errorf("expected %s, got %s", c.expected[j], encounter.Pokemon.Name)
					return
				}
			}
		})
	}
}

func TestGetPokemonData(t *testing.T) {
	cases := []struct {
		name        string
		height      int
		weight      int
		captureRate int
	}{
		{
			name:        "pidgey",
			height:      3,
			weight:      18,
			captureRate: 255,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			client, _ := newTestClient(t)

			pokemon, err := client.GetPokemonData(&c.name)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			// Name, Height and Weight come from /pokemon...
			if pokemon.Name != c.name {
				t.Errorf("expected name %s, got %s", c.name, pokemon.Name)
				return
			}
			if pokemon.Height != c.height {
				t.Errorf("expected height %d, got %d", c.height, pokemon.Height)
				return
			}
			if pokemon.Weight != c.weight {
				t.Errorf("expected weight %d, got %d", c.weight, pokemon.Weight)
				return
			}
			// ...and CaptureRate from /pokemon-species, which proves the merge.
			if pokemon.CaptureRate != c.captureRate {
				t.Errorf("expected capture rate %d, got %d", c.captureRate, pokemon.CaptureRate)
				return
			}
		})
	}
}

// TestCache calls everything twice; with a working cache the server sees each route once.
func TestCache(t *testing.T) {
	client, calls := newTestClient(t)
	location := "canalave-city-area"
	pokemon := "pidgey"

	for range 2 {
		if _, err := client.ListLocations(nil); err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if _, err := client.ListEncounters(&location); err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if _, err := client.GetPokemonData(&pokemon); err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
	}

	routes := []string{
		"/location-area",
		"/location-area/canalave-city-area",
		"/pokemon/pidgey",
		"/pokemon-species/16/",
	}
	for _, route := range routes {
		if calls[route] != 1 {
			t.Errorf("expected 1 request to %s, got %d", route, calls[route])
		}
	}
}

func TestErrors(t *testing.T) {
	unknown := "missingno"
	broken := "broken"

	// Wrapping each call in a func lets one table cover methods with different return types.
	cases := []struct {
		name string
		call func(c *Client) error
	}{
		{
			name: "nil location",
			call: func(c *Client) error {
				_, err := c.ListEncounters(nil)
				return err
			},
		},
		{
			name: "nil pokemon",
			call: func(c *Client) error {
				_, err := c.GetPokemonData(nil)
				return err
			},
		},
		{
			name: "unknown pokemon (404)",
			call: func(c *Client) error {
				_, err := c.GetPokemonData(&unknown)
				return err
			},
		},
		{
			name: "server error (500)",
			call: func(c *Client) error {
				_, err := c.GetPokemonData(&broken)
				return err
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, _ := newTestClient(t)

			if err := c.call(client); err == nil {
				t.Errorf("expected an error")
				return
			}
		})
	}
}
