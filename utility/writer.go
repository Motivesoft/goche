package utility

import (
	"fmt"

	// Internal references
	"goche/status"
)

// Write the engine identification information
func WriteId(engineName string, authorName string) {
	write("id name %s", engineName)
	write("id author %s", authorName)
}

// Write an 'info string'
func WriteInfoString(format string, args ...interface{}) {
	information := fmt.Sprintf(format, args...)
	write("info string %s", information)
}

// Write the registration status
func WriteRegistrationStatus(status status.Status) {
	write("registration %s", status)
}

// Write the 'uciok' message
func WriteUciOk() {
	write("uciok")
}

// Internal print function
func write(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println()
}
