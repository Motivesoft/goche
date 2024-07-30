package uci

import (
	// Internal references
	"goche/identification"
	"goche/logger"
	"goche/status"
	"goche/utility"
)

type Command func(*configuration, string) bool

var commands = map[string]Command{
	"debug":      debugCommand,
	"go":         goCommand,
	"isready":    isreadyCommand,
	"ponderhit":  ponderhitCommand,
	"position":   positionCommand,
	"quit":       quitCommand,
	"register":   registerCommand,
	"setoption":  setoptionCommand,
	"stop":       stopCommand,
	"uci":        uciCommand,
	"ucinewgame": ucinewgameCommand,
}

type configuration struct {
	uciok                bool
	debug                bool
	registrationStatus   status.Status
	copyProtectionStatus status.Status

	// Transient
	registrationWarningIssued bool
}

// NewConfiguration creates a new configuration object with the debug flag set to false.
//
// Parameters:
// - debug: whether to enable debug logging by default
//
// Returns a pointer to the newly created configuration object.
func NewConfiguration() *configuration {
	utility.WriteInfoString("Hello from %s version %s", identification.GetEngineName(), identification.GetVersionName())
	return &configuration{
		uciok:                     false,
		debug:                     false,
		registrationStatus:        status.Checking,
		copyProtectionStatus:      status.Checking,
		registrationWarningIssued: false,
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

	logger.Debug("Received '%s'", command)

	// Check for registration
	if configuration.registrationStatus == status.Error {
		if !configuration.registrationWarningIssued {
			if configuration.uciok {
				if command != "register" {
					configuration.registrationWarningIssued = true
					utility.WriteInfoString("The engine is not registered. Use 'register' to register your engine.")
				}
			}
		}
	}

	// TODO refine this
	if configuration.copyProtectionStatus == status.Error {
		if configuration.uciok {
			logger.Error("The engine copy protection status is not ok")
			utility.WriteInfoString("The engine copy protection status is not ok")

			// Cause the engine to shut down
			// TODO is this the right thing to do?
			return false
		}
	}

	// Call the appropriate command handler
	return commands[command](configuration, arguments)
}

// Process 'debug'
func debugCommand(configuration *configuration, arguments string) bool {
	// Note that this is for verbosity in sending info strings to the caller,
	// not for our own logging, which is handled by the logger package
	configuration.debug = arguments == "on"

	utility.WriteInfoString("Debug mode %s", utility.If(configuration.debug, "enabled", "disabled"))

	return true
}

func goCommand(configuration *configuration, _ string) bool {
	return true
}

func isreadyCommand(configuration *configuration, _ string) bool {
	// TODO check any ongoing activities and wait if necessary

	utility.WriteReadyOk()
	return true
}

func ponderhitCommand(configuration *configuration, _ string) bool {
	return true
}

func positionCommand(configuration *configuration, _ string) bool {
	return true
}

// Process 'quit'
func quitCommand(configuration *configuration, _ string) bool {
	// TODO - check for, and terminate running threads etc

	return false
}

func registerCommand(configuration *configuration, arguments string) bool {

	// TODO create a registration object and do something useful with the name/code used here

	keyword, _ := utility.SplitNextWord(arguments)
	switch keyword {
	case "later":
		configuration.registrationStatus = status.Ok

	case "name":
		configuration.registrationStatus = status.Ok

	case "code":
		configuration.registrationStatus = status.Ok
	}

	// Confirm the change in status
	checkRegistration(configuration)
	return true
}

func setoptionCommand(configuration *configuration, arguments string) bool {
	// Currently no options are supported
	optionName, _ := utility.SplitNextWord(arguments)
	if optionName != "" {
		logger.Warn("Unexpected attempt to configure '%s'", optionName)
	} else {
		logger.Error("Malformed setoption command")
	}

	return true
}

func stopCommand(configuration *configuration, _ string) bool {
	return true
}

// Process 'uci'
func uciCommand(configuration *configuration, _ string) bool {
	if configuration.uciok {
		// Log as error as this is a non-conformance with the UCI spec
		logger.Error("Ignoring duplicate 'uci' command")
	}

	utility.WriteId(identification.GetEngineName(), identification.GetAuthorName())

	// TODO write any options we have

	utility.WriteUciOk()

	// Do the registration stuff
	checkRegistration(configuration)
	checkCopyProtection(configuration)

	configuration.uciok = true

	return true
}

func ucinewgameCommand(configuration *configuration, _ string) bool {
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

	// Reset things
	configuration.registrationWarningIssued = false
}

func checkCopyProtection(configuration *configuration) {
	utility.WriteCopyProtectionStatus(configuration.copyProtectionStatus)

	if configuration.copyProtectionStatus == status.Checking {

		// TODO - check for any copy protection issues and set the status flag accordingly
		configuration.copyProtectionStatus = status.Ok

		// Notify the change in status
		utility.WriteCopyProtectionStatus(configuration.copyProtectionStatus)
	}
}
