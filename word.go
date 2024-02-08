package wikitionaryparser

import (
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Word struct {
	Name      string     `json:"name"`
	Languages []Language `json:"languages"`
}

func parseWord(source string) (Word, error) {
	word := Word{}

	wordName, err := wordNameFromSource(source)
	if err != nil {
		return word, err
	}
	word.Name = wordName

	return word, nil
}

func wordNameFromSource(source string) (string, error) {
	name := ""

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(source))
	if err != nil {
		return name, err
	}

	titleContainer := doc.Find("span.mw-page-title-main")
	if titleContainer.Length() == 0 {
		return name, errors.New("could not find name container")
	}

	name = titleContainer.Text()

	return name, nil
}
