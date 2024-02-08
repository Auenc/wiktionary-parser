package details

import (
	"fmt"
	"os"
	"strings"

	"github.com/auenc/wiktionary-parser/utils"
)

// File mainly here to host testing functions for the rest of the details package
type WordSectionTest struct {
	word     string
	sections map[string]string
}

func LoadWordTestData(word string) (WordSectionTest, error) {
	sectionTest := WordSectionTest{
		word: word,
	}
	testDir := fmt.Sprintf("../testdata/%s", word)
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
