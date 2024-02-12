package wikitionaryparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEtymologyFromSection(t *testing.T) {
	angenSubSectionMap, err := LoadWordTestData("angen/welsh-subsections")
	assert.Nil(t, err)
	tests := []struct {
		name          string
		inputSource   string
		inputLanguage string
		inputWord     string
		output        *Etymology
		expectedError error
	}{
		{
			name:          "can create etymology for welsh word 'angen'",
			inputSource:   angenSubSectionMap.sections["etymology"],
			inputLanguage: "welsh",
			inputWord:     "angen",
			output: &Etymology{
				Language: "welsh",
				Word:     "angen",
				From: &Etymology{
					Language: "middle welsh",
					Word:     "aghen",
					From: &Etymology{
						Language: "proto-brythonic",
						Word:     "*anken",
					},
				},
			},
		},
	}

	for _, test := range tests {
		result, err := etymologyFromSection(test.inputLanguage, test.inputWord, test.inputSource)
		assert.Equal(t, test.expectedError, err)
		assert.Equal(t, test.output, result)
	}
}
