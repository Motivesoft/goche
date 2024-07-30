package uci

import (
	"log"

	// Internal references
	"goche/identification"
	"goche/status"
	"goche/utility"
)

type Command func(*configuration, string) bool

var commands = map[string]Command{
	"debug": debugCommand,
	"quit":  quitCommand,
	"uci":   uciCommand,
}

type configuration struct {
	uciok                bool
	debug                bool
	registrationStatus   status.Status
	copyProtectionStatus status.Status
}

// NewConfiguration creates a new configuration object with the debug flag set to false.
//
// Parameters:
// - debug: whether to enable debug logging by default
//
// Returns a pointer to the newly created configuration object.
func NewConfiguration(debug bool) *configuration {
	utility.WriteInfoString("Hello from %s version %s", identification.GetEngineName(), identification.GetVersionName())
	return &configuration{
		uciok:                false,
		debug:                debug,
		registrationStatus:   status.Checking,
		copyProtectionStatus: status.Checking,
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

	if configuration.registrationStatus == status.Error {
		if configuration.uciok {
			if command != "register" {
				utility.WriteInfoString("The engine is not registered. Use 'register' to register your engine.")

				// TODO decide if we want to be this strict or not
				return true
			}
		}
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
	utility.WriteId(identification.GetEngineName(), identification.GetAuthorName())
	utility.WriteUciOk()

	checkRegistration(configuration)

	return true
}

func registerCommand(configuration *configuration, _ string) bool {
	configuration.registrationStatus = status.Ok
	return true
}

func checkRegistration(configuration *configuration) {
	utility.WriteRegistrationStatus(configuration.registrationStatus)

	if configuration.registrationStatus == status.Checking {

		// TODO - check for registration and reset the status flag
		configuration.registrationStatus = status.Ok

		// Notify the change in status
		utility.WriteRegistrationStatus(configuration.registrationStatus)
	}
}
