package utility

import "fmt"

// Define the Writer interface
type Writer interface {
	WriteUciOk()
}

// Define a concrete type that implements the Writer interface
type ConsoleWriter struct{}

// Implement the Write method for ConsoleWriter
func (cw ConsoleWriter) WriteUciOk() {
	cw.write("uciok")
}

// Implement the Write method for ConsoleWriter
func (cw ConsoleWriter) write(data string) {
	fmt.Println(data)
}
