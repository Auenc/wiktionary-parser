package wikitionaryparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordFromSource(t *testing.T) {
	angenMap, err := LoadWordTestData("angen")
	assert.Nil(t, err)
	tests := []struct {
		name          string
		input         string
		output        string
		expectedError error
	}{
		{
			name:   "can find angen",
			input:  angenMap.sections["source"],
			output: "angen",
		},
	}

	for _, test := range tests {
		result, err := wordNameFromSource(test.input)
		assert.Equal(t, test.expectedError, err)
		assert.Equal(t, test.output, result)
	}
}
