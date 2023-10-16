package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
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
	output := "Usage:\n\n"

	for text, cmd := range generateFunctionMap() {
		output += fmt.Sprintf("\t%s: %s\n", text, cmd.desc)
	}

	return output
}

func commandExit() string {
	return "exiting..."
}

func main() {
	fmt.Println("Welcome to the Pokedex, your one-stop shop for all things Pokemon!")

	cmdChan := make(chan cliCommand)

	// channel for signal processing
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt)

	scanner := bufio.NewScanner(os.Stdin)

	cmds := generateFunctionMap()

	go func() {
		for {
			fmt.Printf("pokedex > ")

			if ok := scanner.Scan(); !ok && scanner.Err() == nil {
				fmt.Printf("received close signal\n")
				cancelChan <- os.Interrupt
				return
			}

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
				cancelChan <- os.Interrupt
				return
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
		case <-cancelChan:
			fmt.Println(cmds["exit"].callback())
			close(cancelChan)
			close(cmdChan)
			fmt.Printf("done closing channels...\n")
			return
		default:
			continue
		}
	}
}
