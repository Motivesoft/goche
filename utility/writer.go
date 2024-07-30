package utility

import "fmt"

// Define the Writer interface
type Writer interface {
	WriteUciOk()
	WriteInfoString(information string)
	WriteId(applicationName string, authorName string)
}

// Define a concrete type that implements the Writer interface
type ConsoleWriter struct{}

// Implement the Write method for ConsoleWriter
func (cw ConsoleWriter) WriteId(applicationName string, authorName string) {
	cw.write("id name " + applicationName)
	cw.write("id author " + authorName)
}

// Implement the Write method for ConsoleWriter
func (cw ConsoleWriter) WriteInfoString(information string) {
	cw.write("info string " + information)
}

// Implement the Write method for ConsoleWriter
func (cw ConsoleWriter) WriteUciOk() {
	cw.write("uciok")
}

// Implement the Write method for ConsoleWriter
func (cw ConsoleWriter) write(data string) {
	fmt.Println(data)
}
