package utils

import (
	"errors"
	"html"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// LoadStringFromFile is just a helper function to load a string from a file. Primarily used for testing
func LoadStringFromFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(strings.TrimSuffix(html.UnescapeString(string(b)), "\n")), nil
}

func StringFromSelector(selector, source string) (string, error) {
	str := ""
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(source))
	if err != nil {
		return str, err
	}
	strContainer := doc.Find("span.mw-page-title-main")
	if strContainer.Length() == 0 {
		return str, errors.New("could not find name container")
	}

	str = strContainer.Text()

	return str, nil
}
