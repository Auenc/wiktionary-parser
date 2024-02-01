package wikitionaryparser

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type wordSectionTest struct {
	word     string
	sections map[string]string
}

func loadWordTestData(word string) (wordSectionTest, error) {
	sectionTest := wordSectionTest{
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
		fileContent, err := loadStringFromFile(fmt.Sprintf("%s/%s", testDir, filePath))
		if err != nil {
			return sectionTest, err
		}
		htmlMap[fileName] = fileContent
	}
	sectionTest.sections = htmlMap

	return sectionTest, nil
}

func TestGetLanguageSection(t *testing.T) {
	angenFull, err := loadWordTestData("angen")
	assert.Nil(t, err)
	tests := []struct {
		name     string
		source   wordSectionTest
		expected []LanguageSection
	}{
		{
			name:   "angen",
			source: angenFull,
			expected: []LanguageSection{
				{
					Name: "balinese",
					Html: angenFull.sections["balinese"],
				},
				{
					Name: "javanese",
					Html: angenFull.sections["javanese"],
				},
				{
					Name: "sundanese",
					Html: angenFull.sections["sundanese"],
				},
				{
					Name: "welsh",
					Html: angenFull.sections["welsh"],
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := getLanguageSections(test.source.sections["source"])
			assert.Nil(t, err)
			for i, section := range test.expected {
				t.Run(section.Name, func(t *testing.T) {
					assert.Equal(t, section.Name, result[i].Name)
					assert.Equal(t, section.Html, result[i].Html)
				})
			}
		})
	}
}
