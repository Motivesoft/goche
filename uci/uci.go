package uci

import "strings"

type Command func(string) bool

var commands = map[string]Command{
	"quit": quitCommand,
}

// ProcessCommand processes a UCI command by extracting the command and arguments
// from the input string and executing the corresponding command function.
//
// Parameters:
// - input: the input string containing the UCI command and arguments.
//
// Returns:
// - bool: true if processing can continue with subsequent commands, false otherwise (e.g. to quit).
func ProcessCommand(input string) bool {
	command, arguments := nextWord(input)

	if command == "" || commands[command] == nil {
		// Illegal commands are silently ignored
		return true
	}

	// Call the appropriate command handler
	return commands[command](arguments)
}

// Process 'quit'
func quitCommand(string) bool {
	return false
}

// Helper function that takes a string and returns the first word and the remaining text
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
