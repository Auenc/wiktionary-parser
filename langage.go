package wikitionaryparser

import "github.com/auenc/wiktionary-parser/details"

type Language struct {
	Name     string            `json:"name"`
	Sections map[string]string `json:"sections"`
	Details  []details.Detail  `json:"details"`
}
