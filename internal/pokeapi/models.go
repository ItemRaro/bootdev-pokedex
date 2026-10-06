package pokeapi

// RespShallowLocations -
type RespShallowLocations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type RespShallowEncounters struct {
	Encounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type PokemonData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	BaseExp     int    `json:"base_experience"`
	CaptureRate int    `json:"capture_rate"`
	IsBaby      bool   `json:"is_baby"`
	IsLegendary bool   `json:"is_legendary"`
	IsMythical  bool   `json:"is_mythical"`
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
