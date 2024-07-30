package uci

import (

	// Internal references
	"fmt"
	"goche/identification"
	"goche/logger"
	"goche/status"
	"goche/utility"
	"strconv"
)

type Command func(*configuration, string) bool

var commands = map[string]Command{
	// Core UCI commands
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

	// Bespoke UCI commands
	"perft": perftCommand,
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

	if command == "" {
		return true
	}

	if commands[command] == nil {
		// Illegal commands are reported and ignored
		logger.Error("Unknown command '%s'", command)

		if configuration.debug {
			utility.WriteInfoString("Unknown command '%s'", command)
		}

		return true
	}

	logger.Debug("Received '%s' command", command)

	if configuration.debug {
		utility.WriteInfoString("Received '%s' command", command)
	}

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
	// TODO implement this

	return true
}

func isreadyCommand(configuration *configuration, _ string) bool {
	// TODO check any ongoing activities and wait if necessary

	utility.WriteReadyOk()
	return true
}

func perftCommand(configuration *configuration, arguments string) bool {
	// TODO implement this

	/*
	   std::cout << "  perft [depth]         - perform a search using a depth and the standard start position" << std::endl;
	   std::cout << "  perft [depth] [fen]   - perform a search using a depth and FEN string" << std::endl;
	   std::cout << "  perft fen [fen]       - perform a search using a FEN string with expected results" << std::endl;
	   std::cout << "  perft file [filename] - perform searches read from a file as FEN strings with expected results" << std::endl;
	*/

	keyword, values := utility.SplitNextWord(arguments)

	if keyword == "" {
		// TODO consider hijacking this to run a pre-determined set of perft tests
		logger.Error("Missing perft command arguments")
		return true
	}

	// Support:
	// - perft [depth]         - perform a search using a depth and the standard start position
	// - perft [depth] [fen]   - perform a search using a depth and FEN string
	if depth, err := strconv.Atoi(keyword); err == nil {
		if values == "" {
			values = FenStartingPosition
		}

		if err = PerftDepth(depth, values); err != nil {
			logger.Error("Error performing perft: %s", err)
		}
	} else if values != "" {
		// Support:
		// - perft fen [fen]       - perform a search using a FEN string containing expected results
		// - perft file [filename] - perform searches read from a file as FEN strings containing expected results
		switch keyword {
		case "fen":
			err = PerftWithFen(values)

		case "file":
			err = PerftWithFile(values)

		default:
			err = fmt.Errorf("unknown command: %s", keyword)
		}

		if err != nil {
			logger.Error("Error performing perft: %s", err)
		}
	} else {
		err = fmt.Errorf("incomplete or invalid command: %s", arguments)
		logger.Error("Error performing perft: %s", err)
	}

	return true
}

func ponderhitCommand(configuration *configuration, _ string) bool {
	// TODO implement this

	return true
}

func positionCommand(configuration *configuration, _ string) bool {
	// TODO implement this

	return true
}

// Process 'quit'
func quitCommand(configuration *configuration, _ string) bool {
	// TODO - check for, and terminate running threads etc, with no sending of bestmove

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

	default:
		configuration.registrationStatus = status.Error
		logger.Error("Malformed register command arguments: %s", arguments)
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

		utility.WriteInfoString("Unknown option '%s'", optionName)
	} else {
		logger.Error("Malformed setoption command")
	}

	return true
}

func stopCommand(configuration *configuration, _ string) bool {
	// TODO see if there is anything ongoing, and terminate it, including causing bestmove to be issued

	return true
}

// Process 'uci'
func uciCommand(configuration *configuration, _ string) bool {
	if configuration.uciok {
		// Log as error as this is a non-conformance with the UCI spec
		logger.Error("'uci' command already issued")

		if configuration.debug {
			utility.WriteInfoString("'uci' command already issued")
		}
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
