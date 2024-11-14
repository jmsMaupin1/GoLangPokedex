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
	for name, _ := range commands {
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

	helpCallback := func(arguments... string) error {
		return printHelp(cliCommandMap)
	}

	exitCallback := func(arguments... string) error {
		shouldExit = true
		return nil
	}
	
	mapCallback := func(arguments... string) error {
		if err := client.GetNextLocationBatch(); err != nil {
			return err
		}

		for _, location := range client.LocationResults.Results {
			fmt.Println(location.Name)
		}

		return nil
	}

	mapbCallback := func(arguments... string) error {
		if err := client.GetPreviousLocationBatch(); err != nil {
			return err
		}

		for _, location := range client.LocationResults.Results {
			fmt.Println(location.Name)
		}

		return nil
	}	

	exploreCallback := func(arguments... string) error {
		if info, err := client.GetLocationInformation(arguments[0]); err != nil {
			return err
		} else {
			for _,pokemon := range info.PokemonEncounters {
				fmt.Println(pokemon.Pokemon.Name)
			}

			return nil
		}
	}

	catchCallback := func(arguments... string) error {
		/*
			chance of capture will be inversely proportional to the base experience
			of the pokemon over the total base experience of all pokemon. So the larger
			the pokemons base experience is, the harder it is to catch.

			currently the total base experience of every pokemon totaled is: 214555
		*/
		const totalPokemonBaseXP = 214555

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
