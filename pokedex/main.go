package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// cliCommand is the default struct for each command within the pokedex
type cliCommand struct {
	name     string
	desc     string
	callback func() string
}

func generateFunctionMap() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:     "help",
			desc:     "Print the usage of the pokedex",
			callback: commandHelp,
		},
		"exit": {
			name:     "exit",
			desc:     "Exit the application cleanly",
			callback: commandExit,
		},
	}
}

func commandHelp() string {
	return "Usage:\nHelp goes here...\n\n"
}

func commandExit() string {
	return "exiting..."
}

func main() {
	fmt.Println("Welcome to the Pokedex, your one-stop shop for all things Pokemon!")

	stopChan := make(chan bool)
	cmdChan := make(chan cliCommand)
	scanner := bufio.NewScanner(os.Stdin)

	cmds := generateFunctionMap()

	go func() {
		for {
			fmt.Printf("pokedex > ")

			_ = scanner.Scan()
			text := strings.TrimSpace(scanner.Text())

			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
				continue
			}

			cmd, ok := cmds[text]
			if !ok {
				fmt.Printf("received garbage command: %s\n", text)
				continue
			}

			if cmd.name == "exit" {
				stopChan <- true
			}

			cmdChan <- cmd
		}
	}()

	for {
		select {
		case cmd := <-cmdChan:
			fmt.Printf("SEBTEST -- we received a command (%s): %s\n", cmd.name, cmd.desc)
			fmt.Println(cmd.callback())
			continue
		case <-stopChan:
			fmt.Println(cmds["exit"].callback())
			close(stopChan)
			close(cmdChan)
			fmt.Printf("done closing channels...\n")
			return
		default:
			continue
		}
	}
}
