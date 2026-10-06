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
		if len(wordList) == 0 {
			continue
		}

		command := wordList[0]
		args := []string{}

		if len(wordList) > 1 {
			args = wordList[1:]
		}

		cmd, ok := cfg.commands[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.callback(cfg, args...); err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
