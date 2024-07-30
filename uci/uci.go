package uci

import (
	"log"

	// Internal references
	"goche/utility"
)

type Command func(*configuration, string) bool

var commands = map[string]Command{
	"debug": debugCommand,
	"quit":  quitCommand,
	"uci":   uciCommand,
}

type Status string

const (
	Ok       Status = "ok"
	Error    Status = "error"
	Checking Status = "checking"
)

type configuration struct {
	debug  bool
	writer utility.Writer
}

// NewConfiguration creates a new configuration object with the debug flag set to false.
//
// Returns a pointer to the newly created configuration object.
func NewConfiguration() *configuration {
	return &configuration{
		debug:  true,
		writer: utility.ConsoleWriter{},
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

// Process 'debug'
func debugCommand(configuration *configuration, arguments string) bool {
	configuration.debug = arguments == "on"
	return true
}

// Process 'quit'
func quitCommand(configuration *configuration, _ string) bool {
	// TODO - check for, and terminate running threads etc

	return false
}

// Process 'uci'
func uciCommand(configuration *configuration, _ string) bool {

	configuration.writer.WriteId(ApplicationName, AuthorName) //configuration.authorName)
	configuration.writer.WriteUciOk()

	return true
}
