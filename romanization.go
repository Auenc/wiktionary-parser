package wikitionaryparser

import (
	"errors"
	"fmt"

	"github.com/auenc/wiktionary-parser/utils"
)

type Romanization struct {
	Of     string `json:"of"`
	OfLink string `json:"ofLink"`
}

func (r Romanization) Name() string {
	return "romanization"
}

func romanizationFromSection(section string) (Romanization, error) {
	rom := Romanization{}

	romSection, err := utils.QueryDocFromstring(section)
	if err != nil {
		return rom, err
	}

	romContainer := romSection.Find("span.form-of-definition-link")
	if romContainer.Length() == 0 {
		return rom, errors.New("could not find romanization container")
	}

	of := romContainer.Text()
	rom.Of = of

	linkTag := romContainer.Find(fmt.Sprintf("a[title='%s']", of))
	if linkTag.Length() == 0 {
		return rom, errors.New("could not find link tag")
	}

	link, exists := linkTag.Attr("href")
	if !exists {
		return rom, errors.New("could not find link href")
	}
	rom.OfLink = link

	return rom, nil
}
