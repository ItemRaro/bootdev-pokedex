package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl(cfg *config) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Println(scanner.Err().Error())
		}

		wordList := cleanInput(scanner.Text())
		cmd, ok := cfg.commands[wordList[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.callback(cfg); err != nil {
			fmt.Println(err)
		}
	}

}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
