package main

import pokecache "github/eldeeishere/pokedexcli/internal"

type locationListResponse struct {
	Results []locationArea `json:"results"`
}

type locationArea struct {
	Name string `json:"name"`
}

type cliCommand struct {
	name        string
	description string
	config      *config
	callback    func(args ...string) error
}

type config struct {
	CurrentOffset int
	cache         *pokecache.Cache
}

type locationAreaResponse struct {
	PokemonEncounters []pokemonEncounter `json:"pokemon_encounters"`
}

type pokemonEncounter struct {
	Pokemon namedAPIResource `json:"pokemon"`
}

type namedAPIResource struct {
	Name string `json:"name"`
}

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Name           string `json:"name"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
	Height int `json:"height"`
}
