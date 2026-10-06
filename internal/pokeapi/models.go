package pokeapi

// RespShallowLocations is a page of location areas.
type RespShallowLocations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// RespShallowEncounters lists the pokemon found in a location area.
type RespShallowEncounters struct {
	Encounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// PokemonData combines fields from the /pokemon and /pokemon-species endpoints.
type PokemonData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	BaseExp     int    `json:"base_experience"`
	CaptureRate int    `json:"capture_rate"` // from /pokemon-species
	IsBaby      bool   `json:"is_baby"`      // from /pokemon-species
	IsLegendary bool   `json:"is_legendary"` // from /pokemon-species
	IsMythical  bool   `json:"is_mythical"`  // from /pokemon-species
	Height      int    `json:"height"`
	Weight      int    `json:"weight"`
	Stats       []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
}
