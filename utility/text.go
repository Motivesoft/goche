package utility

import (
	"strings"
)

// SplitNextWord splits the input string into the first word and the remaining text.
//
// Parameters:
// - input: the input string to be split.
//
// Returns:
// - firstWord: the first word of the input string.
// - remainingText: the remaining text after the first word.
func SplitNextWord(input string) (string, string) {
	input = strings.TrimSpace(input)

	if input == "" {
		return "", ""
	}

	// Find the first space
	spaceIndex := strings.Index(input, " ")

	if spaceIndex == -1 {
		// If there is no space, the entire input is the first word
		return input, ""
	}

	// Split the input into the first word and the remaining text
	firstWord := input[:spaceIndex]
	remainingText := strings.TrimSpace(input[spaceIndex+1:])

	return firstWord, remainingText
}
