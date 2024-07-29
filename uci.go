package main

import "strings"

type Command func(string) bool

var commands = map[string]Command{
	"quit": quitCommand,
}

func quitCommand(string) bool {
	return false
}

func processCommand(input string) bool {
	command, arguments := nextWord(input)

	if command == "" || commands[command] == nil {
		// Illegal commands are silently ignored
		return true
	}

	return commands[command](arguments)
}

func nextWord(input string) (string, string) {
	input = strings.TrimSpace(input)

	if input == "" {
		return "", ""
	}

	// Find the first space
	spaceIndex := strings.Index(input, " ")

	if spaceIndex == -1 {
		// If there is no space, the entire input is the first word
		return input, ""
	}

	// Split the input into the first word and the remaining text
	firstWord := input[:spaceIndex]
	remainingText := strings.TrimSpace(input[spaceIndex+1:])

	return firstWord, remainingText
}
