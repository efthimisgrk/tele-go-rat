package helpers

import (
	"fmt"
	"strings"
)

func ExtractArgument(inputText string) (string, error) {
	
	// Split the input text by spaces
	parts := strings.Split(inputText, " ")

	// Check if there are at least two parts ("/command" and argument)
	if len(parts) >= 2 {
		// Join the remaining parts to get the desired data
		return strings.Join(parts[1:], " "), nil
	}

	// If the input doesn't match the expected format, return an error
	return "", fmt.Errorf("No argument provided")
}
