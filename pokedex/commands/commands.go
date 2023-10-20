package commands

import "fmt"

// CliCommand is the default struct for each command within the pokedex
type CliCommand struct {
	Name     string
	Desc     string
	Callback func() string
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
