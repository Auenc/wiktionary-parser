package wikitionaryparser

import (
	"html"
	"os"
	"strings"
)

// loadStringFromFile is just a helper function to load a string from a file. Primarily used for testing
func loadStringFromFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(strings.TrimSuffix(html.UnescapeString(string(b)), "\n")), nil
}
