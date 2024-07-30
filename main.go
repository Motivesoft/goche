package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	// Internal references
	"goche/uci"
)

func main() {
	// Configure the small number of command line arguments
	inputFile := flag.String("i", "", "filename of UCI commands for testing purposes")
	helpFlag := flag.Bool("h", false, "show this help message and exit")
	versionFlag := flag.Bool("v", false, "print version and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", filepath.Base(os.Args[0]))
		fmt.Println("Options:")
		fmt.Println("  -i filename	" + flag.Lookup("i").Usage)
		fmt.Println("  -h        	" + flag.Lookup("h").Usage)
		fmt.Println("  -v   		" + flag.Lookup("v").Usage)
	}

	flag.Parse()

	// Handle -h
	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	// Handle -v
	if *versionFlag {
		fmt.Printf("%s version %s (%s/%s)\n", uci.GetEngineName(), uci.GetVersionName(), runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	// We expect to take out input from stdin, but allow the user to specify an auto-response input file
	var scanner *bufio.Scanner
	if *inputFile == "" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(*inputFile)
		if err != nil {
			fmt.Println("Error opening input file:", err)
			os.Exit(1)
		}
		defer file.Close()

		// Create a scanner to read from the file
		scanner = bufio.NewScanner(file)
	}

	// Create the environment for the UCI engine
	uciConfiguration := uci.NewConfiguration()

	// Input loop
	for {
		if !scanner.Scan() {
			break
		}

		// Read the input
		input := scanner.Text()

		// Process commands until one of them tells us to break out of loop
		if !uci.ProcessCommand(uciConfiguration, input) {
			break
		}
	}

	// Report operational errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
