package uci

import (
	"goche/utility"
	"log"
)

type Command func(*configuration, string) bool

var commands = map[string]Command{
	"quit": quitCommand,
	"uci":  uciCommand,
}

type configuration struct {
	debug bool
}

// NewConfiguration creates a new configuration object with the debug flag set to false.
//
// Returns a pointer to the newly created configuration object.
func NewConfiguration() *configuration {
	return &configuration{
		debug: true,
	}
}

// ProcessCommand processes a UCI command by extracting the command and arguments
// from the input string and executing the corresponding command function.
//
// Parameters:
// - input: the input string containing the UCI command and arguments.
//
// Returns:
// - bool: true if processing can continue with subsequent commands, false otherwise (e.g. to quit).
func ProcessCommand(configuration *configuration, input string) bool {
	command, arguments := utility.SplitNextWord(input)

	if command == "" || commands[command] == nil {
		// Illegal commands are silently ignored
		return true
	}

	if configuration.debug {
		log.Printf("Received %s", command)
	}

	// Call the appropriate command handler
	return commands[command](configuration, arguments)
}

// Process 'quit'
func quitCommand(configuration *configuration, _ string) bool {
	return false
}

// Process 'uci'
func uciCommand(configuration *configuration, _ string) bool {
	return true
}
