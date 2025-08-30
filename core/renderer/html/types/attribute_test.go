package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAttribute(t *testing.T) {
	// Act
	attr := NewAttribute("name", "value")

	// Assert
	assert.Equal(t, "name", attr.Name())
	assert.Equal(t, "value", attr.Value())
}

func TestAttribute_AddRune(t *testing.T) {
	tests := []struct {
		name          string
		input         rune
		isName        bool
		expectedName  string
		expectedValue string
	}{
		{
			name:          "Add rune to name",
			input:         'a',
			isName:        true,
			expectedName:  "namea",
			expectedValue: "value",
		},
		{
			name:          "Add rune to value",
			input:         'b',
			isName:        false,
			expectedName:  "name",
			expectedValue: "valueb",
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Arrange
			attr := NewAttribute("name", "value")

			// Act
			attr.AddRune(it.input, it.isName)

			// Assert
			assert.Equal(t, it.expectedName, attr.Name())
			assert.Equal(t, it.expectedValue, attr.Value())
		})
	}
}
