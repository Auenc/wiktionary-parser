package wikitionaryparser

import (
	"errors"
	"fmt"
	"html"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/auenc/wiktionary-parser/utils"
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

func indexFromSelection(source string, el *goquery.Selection) (int, error) {
	sectionHtml, err := goquery.OuterHtml(el)
	if err != nil {
		return -1, err
	}

	unescapedNameHtml := html.UnescapeString(sectionHtml)
	index := strings.Index(source, unescapedNameHtml)
	return index, nil
}

// selectionHtmlLength returns the length of the html of the selection
func selectionHtmlLength(source string, el *goquery.Selection) (int, error) {
	sectionHtml, err := goquery.OuterHtml(el)
	if err != nil {
		return -1, err
	}
	sectionHtml = html.UnescapeString(sectionHtml)
	return len(sectionHtml), nil
}

func getSectionsAndTitlesSlice(source, headingTag, selector string) (sections, titles []string, err error) {
	sections = make([]string, 0)
	titles = make([]string, 0)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(source))
	if err != nil {
		return sections, titles, err
	}

	sectionContainer := doc.Find(fmt.Sprintf("%s:has(%s)", headingTag, selector))

	if sectionContainer.Length() == 0 {
		return sections, titles, errors.New("could not find any section containers")
	}

	sectionNames := make([]string, 0)
	sectionIndexes := make([]int, 0)

	sectionContainer.EachWithBreak(func(i int, el *goquery.Selection) bool {
		sectionName := el.Find(selector).Text()
		sectionNames = append(sectionNames, sectionName)

		// creating variable this way so we don't scope err to the annonymous function
		var languageStart int
		languageStart, err = indexFromSelection(source, el)
		if languageStart == -1 {
			err = fmt.Errorf("could not find starting index of section: %s", sectionName)
			return false
		}
		sectionIndexes = append(sectionIndexes, languageStart)
		return true
	})
	if err != nil {
		return sections, titles, err
	}

	contentEnd := len(source)
	for i, start := range sectionIndexes {
		end := contentEnd
		if i != len(sectionIndexes)-1 {
			end = sectionIndexes[i+1]
		}
		html := source[start:end]
		titles = append(titles, strings.ToLower(sectionNames[i]))
		sections = append(sections, strings.TrimSuffix(html, "\n"))
	}

	return sections, titles, err
}

func getSectionMap(source string, headingTag, selector string) (map[string]string, error) {
	sectionsMap := make(map[string]string)

	sections, titles, err := getSectionsAndTitlesSlice(source, headingTag, selector)
	if err != nil {
		return sectionsMap, err
	}

	for i, section := range sections {
		sectionsMap[titles[i]] = section
	}

	return sectionsMap, nil
}

func extractSection(source *string, selector, containerSelector string) (string, error) {
	doc, err := utils.QueryDocFromstring(*source)
	if err != nil {
		return "", err
	}

	if containerSelector == "" {
		containerSelector = "body"
	}

	container := doc.Find(containerSelector)
	containerHtml, err := container.Html()
	if err != nil {
		return "", nil
	}
	containerHtml = html.UnescapeString(containerHtml)

	if len(containerHtml) == 0 {
		return "", nil
	}

	element := container.Find(selector)

	first := element.First()
	if first.Length() == 0 {
		return "", nil
	}

	extractStart, err := indexFromSelection(containerHtml, first)
	if err != nil {
		return "", err
	}
	if extractStart == -1 {
		return "", nil
	}

	extractLength, err := selectionHtmlLength(containerHtml, first)
	if err != nil {
		return "", err
	}

	extractEnd := extractStart + extractLength
	if extractEnd > len(containerHtml) {
		extractEnd = len(containerHtml) - 1
	}
	extracted := containerHtml[0:extractEnd]
	*source = strings.Replace(*source, extracted, "", 1)

	return extracted, nil
}

func splitBySelection(source, selector, containerSelector string) ([]string, error) {
	sections := make([]string, 0)
	replaced := source
	extracted, err := extractSection(&replaced, selector, containerSelector)
	if err != nil {
		return sections, err
	}
	for extracted != "" {
		sections = append(sections, extracted)
		extracted, err = extractSection(&replaced, selector, containerSelector)
		if err != nil {
			return sections, err
		}
	}

	return sections, nil
}
