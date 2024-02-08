package wikitionaryparser

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/auenc/wiktionary-parser/utils"
	"github.com/stretchr/testify/assert"
)

type WordSectionTest struct {
	word     string
	sections map[string]string
}

func LoadWordTestData(word string) (WordSectionTest, error) {
	sectionTest := WordSectionTest{
		word: word,
	}
	testDir := fmt.Sprintf("testdata/%s", word)
	files, err := os.ReadDir(testDir)
	if err != nil {
		return sectionTest, err
	}

	htmlMap := make(map[string]string)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := file.Name()
		fileNameWithExt := strings.ReplaceAll(filePath, testDir, "")
		fileName := strings.ReplaceAll(fileNameWithExt, ".html", "")
		fileContent, err := utils.LoadStringFromFile(fmt.Sprintf("%s/%s", testDir, filePath))
		if err != nil {
			return sectionTest, err
		}
		htmlMap[fileName] = fileContent
	}
	sectionTest.sections = htmlMap

	return sectionTest, nil
}
func TestGetSectionMap(t *testing.T) {
	angenSubSectionMap, err := LoadWordTestData("angen/welsh-subsections")
	assert.Nil(t, err)
	angenMap, err := LoadWordTestData("angen")
	assert.Nil(t, err)

	angenExpectedSections := make(map[string]string)
	for sectionName, sectionData := range angenMap.sections {
		if sectionName == "section-container" || sectionName == "source" {
			continue
		}
		angenExpectedSections[sectionName] = sectionData
	}

	angenExpectedSubSections := make(map[string]string)
	for subName, subData := range angenSubSectionMap.sections {
		angenExpectedSubSections[subName] = subData
	}

	tests := []struct {
		name     string
		input    string
		tag      string
		selector string
		expected map[string]string
	}{
		{
			name:     "basic section contaienr",
			input:    strings.TrimSpace("<h2><span>section 1</span></h2><p>text for section 1</p><h2><span>section 2</span></h2><p>section two text</p>"),
			tag:      "h2",
			selector: "span",
			expected: map[string]string{"section 1": "<h2><span>section 1</span></h2><p>text for section 1</p>", "section 2": "<h2><span>section 2</span></h2><p>section two text</p>"},
		},
		{
			name:     "can pull out language sections for word: angen",
			tag:      "h2",
			selector: "span.mw-headline",
			input:    angenMap.sections["section-container"],
			expected: angenExpectedSections,
		},
		{
			name:     "can pull the welsh language subsections",
			tag:      "h3",
			selector: "span.mw-headline",
			input:    angenMap.sections["welsh"],
			expected: angenExpectedSubSections,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := getSectionMap(test.input, test.tag, test.selector)
			assert.Nil(t, err)
			assert.Equal(t, len(test.expected), len(result))
			for expectedSection, expectedContent := range test.expected {
				assert.Equal(t, expectedContent, result[expectedSection])
			}
		})
	}
}
