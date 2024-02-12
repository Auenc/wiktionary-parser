package wikitionaryparser

import (
	"strings"

	"github.com/auenc/wiktionary-parser/utils"
)

type Etymology struct {
	Language string     `json:"language"`
	Word     string     `json:"word"`
	From     *Etymology `json:"etymology"`
}

func etymologyFromSection(language, word, section string) (*Etymology, error) {
	et := &Etymology{Language: language, Word: word}

	ets, err := splitBySelection(section, "i.mention", "p")
	if err != nil {
		return et, err
	}
	currentEt := et
	for _, etSection := range ets {
		etSection, err = utils.RemoveNodesFromString(etSection, "span.annotation-paren")
		if err != nil {
			return et, err
		}
		etSection, err = utils.RemoveNodesFromString(etSection, "span.mention-gloss-double-quote")
		if err != nil {
			return et, err
		}
		etSection, err = utils.RemoveNodesFromString(etSection, "span.mention-gloss")
		if err != nil {
			return et, err
		}
		etSection = strings.Replace(etSection, ")", "", 1)
		etSection = strings.Replace(etSection, ",", "", 1)
		isFrom := isFromEty(etSection, strings.ToLower(word) == "from")

		// Making the assumption that wiktionary will always have the "from" etymology first. Stop when we hit first non-from etymology
		// We should expand upon this to allow us to capture a structure like from(non-from->from->from)->from
		if !isFrom {
			return et, nil
		}
		fromEt, err := etymologyFromFromEt(etSection)
		if err != nil {
			return et, err
		}
		currentEt.From = &fromEt
		currentEt = currentEt.From
	}

	return et, nil
}

func etymologyFromFromEt(source string) (Etymology, error) {
	et := Etymology{}
	languageName, err := utils.StringFromSelector("a.extiw", source)
	if err != nil {
		return et, err
	}
	languageWord, err := utils.StringFromSelector("i.mention", source)
	if err != nil {
		return et, err
	}
	et.Language = strings.ToLower(languageName)
	et.Word = strings.ToLower(languageWord)

	return et, nil
}

func isFromEty(source string, isWordFrom bool) bool {
	fromCountRequired := 1
	if isWordFrom {
		fromCountRequired = 2
	}
	return strings.Count(strings.ToLower(source), "from") == fromCountRequired
}
