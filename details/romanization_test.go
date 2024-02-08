package details

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRomanizationFromSection(t *testing.T) {
	angenMap, err := LoadWordTestData("angen")
	assert.Nil(t, err)
	tests := []struct {
		name          string
		input         string
		output        Romanization
		expectedError error
	}{
		{
			name:  "can create balinese romanization",
			input: angenMap.sections["balinese"],
			output: Romanization{
				Of:     "ᬳᬗᭂᬦ᭄",
				OfLink: "/wiki/%E1%AC%B3%E1%AC%97%E1%AD%82%E1%AC%A6%E1%AD%84#Balinese",
			},
		},
	}

	for _, test := range tests {
		fmt.Println("hellp", test.input)
		result, err := romanizationFromSection(test.input)
		assert.Equal(t, test.expectedError, err)
		assert.Equal(t, test.output, result)
	}
}
