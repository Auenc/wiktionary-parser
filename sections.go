package wikitionaryparser

import (
	"fmt"
	"html"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type LanguageSection struct {
	Name string
	Html string
}

func getLanguageSections(pageSource string) ([]LanguageSection, error) {
	languageSections := make([]LanguageSection, 0)

	languageNames := make([]string, 0)
	languageSectionIndexes := make([]int, 0)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(pageSource))
	if err != nil {
		return languageSections, err
	}

	contentArea := doc.Find("div.mw-content-ltr.mw-parser-output")
	sectionTitles := contentArea.Find("h2:has(span.mw-headline)")

	contentAreaHtml, err := contentArea.Html()
	if err != nil {
		return languageSections, err
	}
	sectionTitles.EachWithBreak(func(i int, el *goquery.Selection) bool {
		languageName := el.Find("span.mw-headline").Text()
		languageNames = append(languageNames, languageName)

		// creating variable this way so we don't scope err to the annonymous function
		var languageNameHtml string
		languageNameHtml, err = el.Html()
		if err != nil {
			return false
		}
		languageNameHtml = strings.TrimSpace(languageNameHtml)
		languageStart := strings.Index(contentAreaHtml, fmt.Sprintf("<h2>%s</h2>", languageNameHtml))
		if languageStart == -1 {
			err = fmt.Errorf("could not find starting index of language section: %s", languageName)
			return false
		}
		languageSectionIndexes = append(languageSectionIndexes, languageStart)
		return true
	})
	if err != nil {
		return languageSections, err
	}

	contentEnd := len(contentAreaHtml)
	for i, start := range languageSectionIndexes {
		end := contentEnd
		if i != len(languageSectionIndexes)-1 {
			end = languageSectionIndexes[i+1]
		}
		html := html.UnescapeString(contentAreaHtml[start:end])
		section := LanguageSection{
			Name: strings.ToLower(languageNames[i]),
			Html: strings.TrimSuffix(html, "\n"),
		}

		languageSections = append(languageSections, section)
	}

	return languageSections, nil
}
