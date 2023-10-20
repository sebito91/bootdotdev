package commands

import (
	"fmt"
	"net/url"
)

// CliCommand is the default struct for each command within the pokedex
type CliCommand struct {
	Name     string
	Desc     string
	Callback func() string
}

// Config is a helper-struct that provides additional context to the CliCommand interaction
type Config struct {
	Next     *url.URL
	Previous *url.URL
}

func GetConfig(startURL string) (*Config, error) {
	if startURL == "" {
		fmt.Println("WARN  -- no `startURL` given, using default...")
		startURL = "https://pokeapi.co/api/v2/location-area/"
	}

	u, err := url.Parse(startURL)
	if err != nil {
		return nil, err
	}

	return &Config{
		Next:     u,
		Previous: u,
	}, nil
}

func GenerateFunctionMap() map[string]CliCommand {
	return map[string]CliCommand{
		"help": {
			Name:     "help",
			Desc:     "Print the usage of the pokedex",
			Callback: commandHelp,
		},
		"map": {
			Name:     "map",
			Desc:     "iterate forward through the set of map locations",
			Callback: commandMap,
		},
		"mapb": {
			Name:     "mapb",
			Desc:     "iterater backward through the set of map locations",
			Callback: commandMapb,
		},
		"exit": {
			Name:     "exit",
			Desc:     "Exit the application cleanly",
			Callback: commandExit,
		},
	}
}

func commandHelp() string {
	output := "Usage:\n\n"

	for text, cmd := range GenerateFunctionMap() {
		output += fmt.Sprintf("\t%s: %s\n", text, cmd.Desc)
	}

	return output
}

func commandExit() string {
	return "exiting..."
}

func commandMap() string {
	return "doing some stuff moving forward"
}

func commandMapb() string {
	return "doing some stuff moving backward"
}
