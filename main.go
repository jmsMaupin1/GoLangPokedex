package main

import (
	"fmt"
	"bufio"
	"os"
)

type cliCommand struct {
	name string
	desc string
	callback func() error
}

const prefix = "pokedex > "

func (c cliCommand) getHelpString() string {
	return fmt.Sprintf("%s: %s", c.name, c.desc)
}

func printHelp(commands map[string]cliCommand) error {
	for name, _ := range commands {
		fmt.Println(commands[name].getHelpString())
	}

	return nil
}

func main() {
	var cliCommandMap map[string]cliCommand
	var shouldExit bool

	helpCallback := func() error {
		return printHelp(cliCommandMap)
	}

	exitCallback := func() error {
		shouldExit = true
		return nil
	}

	cliCommandMap = map[string]cliCommand {
		"help": {
			name: "help",
			desc: "Displays this help message",
			callback: helpCallback,
		},
		"exit": {
			name: "exit",
			desc: "Exits the program",
			callback: exitCallback,
		},
	}


	cliCommandMap["help"].callback()
	scanner := bufio.NewScanner(os.Stdin)

	for ;!shouldExit;{
		fmt.Print(prefix)
		scanner.Scan()
		
		if command, ok := cliCommandMap[scanner.Text()]; !ok {
			helpCallback()
		} else {
			command.callback()
		}
	}	
}
