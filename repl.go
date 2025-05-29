package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	pokecache "github/eldeeishere/pokedexcli/internal"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func startRepl() {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		formatedText := cleanInput(input)
		if scanner.Err() != nil {
			fmt.Println("Error reading input:", scanner.Err())
			continue
		}
		if len(formatedText) == 0 {
			fmt.Println("No input provided, please try again.")
			continue
		}
		if cmd, ok := commands[formatedText[0]]; ok {
			err := cmd.callback(formatedText[1:]...)
			if err != nil {
				fmt.Printf("Error executing command '%s': %v\n", formatedText[0], err)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", formatedText[0])
		}

	}
}

func cleanInput(text string) []string {
	toLower := strings.ToLower(text)
	split := strings.Fields(toLower)

	return split
}

func apiCallGet(endpoint string) (*http.Response, error) {
	res, err := http.Get(apiURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("error fetching data from the API: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("received status code %d from the API", res.StatusCode)
	}

	return res, nil
}

func fetchOrCache(endpoint string, target any, cache *pokecache.Cache) error {
	if data, ok := cache.Get(endpoint); ok {
		return json.Unmarshal(data, target)
	}

	res, err := apiCallGet(endpoint)
	if err != nil {
		return fmt.Errorf("error fetching data from the API: %w", err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	cache.Add(endpoint, data)

	return json.Unmarshal(data, target)

}

func randomChance(baseExperience int) bool {
	rand.Seed(time.Now().UnixNano())
	chance := 100 - (baseExperience / 2)
	randomValue := rand.Intn(100) + 1 // Random value between 1 and 100
	if chance < 5 {
		chance = 5 // Ensure a minimum chance of 5%
	}
	if chance > 95 {
		chance = 95 // Ensure a maximum chance of 95%
	}
	return randomValue <= chance
}
