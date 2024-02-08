package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

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
		fileContent, err := LoadStringFromFile(fmt.Sprintf("%s/%s", testDir, filePath))
		if err != nil {
			return sectionTest, err
		}
		htmlMap[fileName] = fileContent
	}
	sectionTest.sections = htmlMap

	return sectionTest, nil
}

func TeststringFromSelector(t *testing.T) {
	angenMap, err := LoadWordTestData("angen")
	assert.Nil(t, err)
	tests := []struct {
		name          string
		selector      string
		input         string
		output        string
		expectedError error
	}{
		{
			name:          "empty string",
			selector:      "",
			input:         "",
			output:        "",
			expectedError: errors.New("could not find name container"),
		},
		{
			name:     "can find angen",
			selector: "span.mw-page-title-main",
			input:    angenMap.sections["source"],
			output:   "angen",
		},
	}

	for _, test := range tests {
		result, err := StringFromSelector(test.selector, test.input)
		assert.Equal(t, test.expectedError, err)
		assert.Equal(t, test.output, result)
	}
}
