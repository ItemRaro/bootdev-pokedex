package pokeapi

type locationResults struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// RespShallowLocations -
type RespShallowLocations struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []locationResults `json:"results"`
}
