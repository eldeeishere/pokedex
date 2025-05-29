package main

import (
	"fmt"
	pokecache "github/eldeeishere/pokedexcli/internal"
	"os"
	"strings"
	"time"
)

const apiURL = "https://pokeapi.co/api/v2/"

var commands map[string]cliCommand

var pokemonInBag map[string]Pokemon

func commandExit(args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandCatch(args ...string) error {
	if len(args) < 1 {
		fmt.Println("Please provide a Pokémon name to catch.")
		return nil
	}
	pokemonName := strings.ToLower(args[0])
	endpoint := "pokemon/" + pokemonName + "/"
	pokemon := Pokemon{}
	if err := fetchOrCache(endpoint, &pokemon, commands["catch"].config.cache); err != nil {
		fmt.Printf("Error fetching Pokémon data: %v\n", err)
		return nil
	}
	if _, exists := pokemonInBag[pokemon.Name]; exists {
		fmt.Printf("You already have %s in your bag!\n", pokemon.Name)
		return nil
	}
	cachedPokemon := randomChance(pokemon.BaseExperience)
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	if cachedPokemon {
		fmt.Printf("%s was caught!!\n", pokemon.Name)
		pokemonInBag[pokemon.Name] = pokemon
		return nil
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}
}

func commandInspect(args ...string) error {
	if len(args) < 1 {
		fmt.Println("Please provide a Pokémon name to inspect.")
		return nil
	}
	pokemonName := strings.ToLower(args[0])
	pokemon, exists := pokemonInBag[pokemonName]
	if !exists {
		fmt.Printf("You don't have a Pokémon named %s in your bag.\n", pokemonName)
		return nil
	}

	fmt.Printf("Inspecting %s...\n", pokemon.Name)
	fmt.Printf("Base Experience: %d\n", pokemon.BaseExperience)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("- %s: %d \n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("- %s\n", t.Type.Name)
	}

	return nil
}

func commandPokebag(args ...string) error {
	if len(pokemonInBag) == 0 {
		fmt.Println("Your Pokebag is empty.")
		return nil
	}
	fmt.Println("Your Pokebag contains:")
	for name, pokemon := range pokemonInBag {
		fmt.Printf("- %s (Base Experience: %d)\n", name, pokemon.BaseExperience)
	}
	return nil
}

func commandExplore(args ...string) error {
	if len(args) < 1 {
		fmt.Println("Please provide a location area name to explore.")
		return nil
	}
	page := commands["explore"].config
	area := locationAreaResponse{}
	locationName := strings.ToLower(args[0])
	endpoint := "location-area/" + locationName + "/"

	if err := fetchOrCache(endpoint, &area, page.cache); err != nil {
		fmt.Printf("Error fetching location data: %v\n", err)
		return err
	}

	if len(area.PokemonEncounters) == 0 {
		fmt.Println("No Pokémon found in this location area.")
		return nil
	}
	fmt.Printf("Exploring %s...\n", locationName)
	fmt.Println("Found Pokemon:")
	for _, location := range area.PokemonEncounters {
		fmt.Printf("- %s\n", location.Pokemon.Name)

	}

	return nil
}

func commandMap(args ...string) error {
	page := commands["map"].config
	m := locationListResponse{}
	// Ideme na ďalšiu stránku
	page.CurrentOffset += 20
	endpoint := "location-area/?offset=" + fmt.Sprint(page.CurrentOffset) + "&limit=20"
	if err := fetchOrCache(endpoint, &m, page.cache); err != nil {
		fmt.Printf("Error fetching location data: %v\n", err)
		page.CurrentOffset -= 20 // Vrátime späť ak sa nepodarilo
		return err
	}
	for _, location := range m.Results {
		fmt.Println(location.Name)

	}
	return nil
}

func commandMapb(args ...string) error {
	page := commands["map"].config

	// Kontrola či môžeme ísť späť
	if page.CurrentOffset-20 < 0 {
		fmt.Println("You are already on the first page, there is no previous page.")
		return nil
	}

	// Ideme na predchádzajúcu stránku
	page.CurrentOffset -= 20

	m := locationListResponse{}
	endpoint := "location-area/?offset=" + fmt.Sprint(page.CurrentOffset) + "&limit=20"
	if err := fetchOrCache(endpoint, &m, page.cache); err != nil {
		fmt.Printf("Error fetching location data: %v\n", err)
		page.CurrentOffset += 20 // Vrátime späť ak sa nepodarilo
		return err
	}

	for _, location := range m.Results {
		fmt.Println(location.Name)
	}
	return nil

}

func init() {
	pokemonInBag = make(map[string]Pokemon)

	sharedConfig := &config{
		CurrentOffset: -20,
		cache:         pokecache.NewCache(5 * time.Second),
	}

	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas",
			config:      sharedConfig,
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			config:      sharedConfig,
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays pokemons in the current location area",
			config:      sharedConfig,
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catches a pokemon",
			config:      sharedConfig,
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a stats of a pokemon that you caught",
			config:      sharedConfig,
			callback:    commandInspect,
		},
		"pokebag": {
			name:        "pokebag",
			description: "Inspects a pokemon in your bag",
			config:      sharedConfig,
			callback:    commandPokebag,
		},
	}
}
