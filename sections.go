package wikitionaryparser

import (
	"fmt"
	"html"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type LanguageSection struct {
	Name        string
	Html        string
	Subsections map[string]string
}

func getLanguageSections(pageSource string) ([]LanguageSection, error) {
	languageSections := make([]LanguageSection, 0)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(pageSource))
	if err != nil {
		return languageSections, err
	}
	contentArea := doc.Find("div.mw-content-ltr.mw-parser-output")

	contentAreaHtml, err := contentArea.Html()
	if err != nil {
		return languageSections, err
	}

	sectionMap, err := getSectionMap(html.UnescapeString(contentAreaHtml), "h2", "span.mw-headline")
	if err != nil {
		return languageSections, err
	}

	for name, sectionHtml := range sectionMap {
		subsectionMap, err := getSectionMap(html.UnescapeString(sectionHtml), "h3", "span.mw-headline")
		if err != nil {
			return languageSections, err
		}
		section := LanguageSection{
			Name:        name,
			Html:        sectionHtml,
			Subsections: subsectionMap,
		}
		languageSections = append(languageSections, section)
	}

	return languageSections, nil
}

func getSectionMap(source string, headingTag, selector string) (map[string]string, error) {
	sections := make(map[string]string)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(source))
	if err != nil {
		return sections, err
	}

	sectionContainer := doc.Find(fmt.Sprintf("%s:has(%s)", headingTag, selector))

	sectionNames := make([]string, 0)
	sectionIndexes := make([]int, 0)

	sectionContainer.EachWithBreak(func(i int, el *goquery.Selection) bool {
		sectionName := el.Find(selector).Text()
		sectionNames = append(sectionNames, sectionName)

		// creating variable this way so we don't scope err to the annonymous function
		var sectionNameHtml string
		sectionNameHtml, err = el.Html()
		if err != nil {
			return false
		}
		unescapedNameHtml := html.UnescapeString(sectionNameHtml)
		languageStart := strings.Index(source, fmt.Sprintf("<%s>%s</%s>", headingTag, unescapedNameHtml, headingTag))
		if languageStart == -1 {
			err = fmt.Errorf("could not find starting index of section: %s\nsearching: %s", sectionName, unescapedNameHtml)
			return false
		}
		sectionIndexes = append(sectionIndexes, languageStart)
		return true
	})
	if err != nil {
		return sections, err
	}

	contentEnd := len(source)
	for i, start := range sectionIndexes {
		end := contentEnd
		if i != len(sectionIndexes)-1 {
			end = sectionIndexes[i+1]
		}
		html := source[start:end]
		sections[strings.ToLower(sectionNames[i])] = strings.TrimSuffix(html, "\n")
	}

	return sections, nil
}
