package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/sebito91/bootdotdev/pokedex/commands"
)

func main() {
	fmt.Println("Welcome to the Pokedex, your one-stop shop for all things Pokemon!")

	cmdChan := make(chan commands.CliCommand)

	// channel for signal processing
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt)

	scanner := bufio.NewScanner(os.Stdin)

	cmds := commands.GenerateFunctionMap()

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

			if cmd.Name == "exit" {
				cancelChan <- os.Interrupt
				return
			}

			cmdChan <- cmd
		}
	}()

	for {
		select {
		case cmd := <-cmdChan:
			fmt.Printf("SEBTEST -- we received a command (%s): %s\n", cmd.Name, cmd.Desc)
			fmt.Println(cmd.Callback())
			continue
		case <-cancelChan:
			fmt.Println(cmds["exit"].Callback())
			close(cancelChan)
			close(cmdChan)
			fmt.Printf("done closing channels...\n")
			return
		default:
			continue
		}
	}
}
