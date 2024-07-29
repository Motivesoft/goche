package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
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

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Println("Version 0.0.1")
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

	fmt.Println("Enter text (type 'quit' to exit):")

	for {
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		if strings.TrimSpace(input) == "quit" {
			fmt.Println("Exiting...")
			break
		}

		fmt.Println("You entered:", input)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
