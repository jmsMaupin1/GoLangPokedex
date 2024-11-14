package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"math/rand"
	"github.com/jmsMaupin1/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name string
	desc string
	callback func(arguments... string) error
}

const prefix = "pokedex > "

func (c cliCommand) getHelpString() string {
	return fmt.Sprintf("%s: %s", c.name, c.desc)
}

func printHelp(commands map[string]cliCommand) error {
	for name := range commands {
		fmt.Println(commands[name].getHelpString())
	}

	return nil
}

func main() {
	var cliCommandMap map[string]cliCommand
	var shouldExit bool

	client, err := pokeapi.NewClient("https://pokeapi.co/api/v2")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error initializing client: %v", err))
	}

	helpCallback := func(arguments ...string) error {
		return printHelp(cliCommandMap)
	}

	exitCallback := func(arguments ...string) error {
		shouldExit = true
		return nil
	}
	
	mapCallback := func(arguments ...string) error {
		if err := client.GetNextLocationBatch(); err != nil {
			return err
		}

		for _, location := range client.LocationResults.Results {
			fmt.Println(location.Name)
		}

		return nil
	}

	mapbCallback := func(arguments ...string) error {
		if err := client.GetPreviousLocationBatch(); err != nil {
			return err
		}

		for _, location := range client.LocationResults.Results {
			fmt.Println(location.Name)
		}

		return nil
	}	

	exploreCallback := func(arguments ...string) error {
		if info, err := client.GetLocationInformation(arguments[0]); err != nil {
			return err
		} else {
			for _,pokemon := range info.PokemonEncounters {
				fmt.Println(pokemon.Pokemon.Name)
			}

			return nil
		}
	}

	inspect := func(arguments ...string) error {
		pokemon, ok := client.CaughtPokemon[arguments[0]]
		if !ok {
			return fmt.Errorf("You havent caught %s yet", arguments[0])
		}

		fmt.Println(fmt.Sprintf("Name %s", pokemon.Name))
		fmt.Println(fmt.Sprintf("Height: %d", pokemon.Height))
		fmt.Println(fmt.Sprintf("Weight: %d", pokemon.Weight))
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Println(fmt.Sprintf(" - %s: %d", stat.Stat.Name, stat.BaseStat))
		}
		fmt.Println("Types:")
		for _, t := range pokemon.Types {
			fmt.Println(fmt.Sprintf(" - %s", t.Type.Name))
		}
		
		return nil
	}

	catchCallback := func(arguments... string) error {	
		if pokemon, err := client.GetPokemonInformation(arguments[0]); err != nil {
			return err
		} else {
			if rand.Intn(pokemon.BaseExperience) > 40 {
				fmt.Println(fmt.Sprintf("You caught %s!", pokemon.Name))
				fmt.Println(fmt.Sprintf("Adding %s to your pokedex", pokemon.Name))
				client.CaughtPokemon[pokemon.Name] = pokemon
			} else {
				fmt.Println(fmt.Sprintf("%s got away :(", pokemon.Name))
			}
			return nil
		}
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
		"map": {
			name: "map",
			desc: "lists the next 20 locations",
			callback: mapCallback,
		},
		"mapb": {
			name: "mapb",
			desc: "lists the previous 20 locations",
			callback: mapbCallback,
		},
		"explore": {
			name: "explore <area name>",
			desc: "Lists all pokemon in a given area",
			callback: exploreCallback,
		},
		"catch": {
			name: "catch <pokemon name>",
			desc: "Attempts to catch a given pokemon name",
			callback: catchCallback,
		},
		"inspect": {
			name: "inspect <pokemon name>",
			desc: "Inspects a given pokemon",
			callback: inspect,
		},
	}


	cliCommandMap["help"].callback()
	scanner := bufio.NewScanner(os.Stdin)

	for ;!shouldExit;{
		fmt.Print(prefix)
		scanner.Scan()

		commandAndArgs := strings.Split(scanner.Text(), " ")
		
		if command, ok := cliCommandMap[commandAndArgs[0]]; !ok {
			helpCallback()
		} else {
			err := command.callback(commandAndArgs[1:]...)
			if err != nil {
				fmt.Println(fmt.Sprintf("Error: %v", err))
			}
		}
	}	
}
