package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	// Internal references
	"goche/identification"
	"goche/uci"
)

func main() {
	// Configure the small number of command line arguments
	inputFile := flag.String("i", "", "filename of UCI commands for testing purposes")
	logFile := flag.String("l", "", "filename for logging output")
	debugFlag := flag.Bool("d", false, "enable debug logging")
	helpFlag := flag.Bool("h", false, "show this help message and exit")
	versionFlag := flag.Bool("v", false, "print version and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", filepath.Base(os.Args[0]))
		fmt.Println("Options:")
		fmt.Println("  -i filename	" + flag.Lookup("i").Usage)
		fmt.Println("  -l filename	" + flag.Lookup("l").Usage)
		fmt.Println("  -d   		" + flag.Lookup("d").Usage)
		fmt.Println("  -v   		" + flag.Lookup("v").Usage)
		fmt.Println("  -h        	" + flag.Lookup("h").Usage)
	}

	flag.Parse()

	// Handle -h
	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	// Handle -v
	if *versionFlag {
		fmt.Printf("%s version %s (%s/%s)\n", identification.GetEngineName(), identification.GetVersionName(), runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	// Log to stderr or the file specified
	if *logFile == "" {
		log.SetOutput(os.Stderr)
	} else {
		// Optionally, add 'os.O_APPEND|' to open the file in append mode
		logFile, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening log file:", err)
			os.Exit(1)
		}
		defer logFile.Close()

		log.SetOutput(logFile)
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
	uciConfiguration := uci.NewConfiguration(*debugFlag)

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
