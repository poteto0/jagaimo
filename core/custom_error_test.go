package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		cause    error
		expected string
	}{
		{
			name:     "NewNetworkError with cause",
			message:  "Network error",
			cause:    errors.New("connection refused"),
			expected: "Network error: connection refused",
		},
		{
			name:     "NewNetworkError without cause",
			message:  "Network error",
			cause:    nil,
			expected: "Network error",
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Arrange
			err := NewNetworkError(it.message, it.cause)

			// Act & Assert
			assert.Equal(t, it.expected, err.Error())
		})
	}
}

func TestUnexpectedInputError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		cause    error
		expected string
	}{
		{
			name:     "NewUnexpectedInputError with cause",
			message:  "Unexpected input",
			cause:    errors.New("invalid input"),
			expected: "Unexpected input: invalid input",
		},
		{
			name:     "NewUnexpectedInputError without cause",
			message:  "Unexpected input",
			cause:    nil,
			expected: "Unexpected input",
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Arrange
			err := NewUnexpectedInputError(it.message, it.cause)

			// Act & Assert
			assert.Equal(t, it.expected, err.Error())
		})
	}
}

func TestInvalidUIError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		cause    error
		expected string
	}{
		{
			name:     "NewUnexpectedInputError with cause",
			message:  "Unexpected input",
			cause:    errors.New("invalid input"),
			expected: "Unexpected input: invalid input",
		},
		{
			name:     "NewUnexpectedInputError without cause",
			message:  "Unexpected input",
			cause:    nil,
			expected: "Unexpected input",
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Arrange
			err := NewInvalidUIError(it.message, it.cause)

			// Act & Assert
			assert.Equal(t, it.expected, err.Error())
		})
	}
}

func TestOtherError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		cause    error
		expected string
	}{
		{
			name:     "NewUnexpectedInputError with cause",
			message:  "Unexpected input",
			cause:    errors.New("invalid input"),
			expected: "Unexpected input: invalid input",
		},
		{
			name:     "NewUnexpectedInputError without cause",
			message:  "Unexpected input",
			cause:    nil,
			expected: "Unexpected input",
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Arrange
			err := NewOtherError(it.message, it.cause)

			// Act & Assert
			assert.Equal(t, it.expected, err.Error())
		})
	}
}
