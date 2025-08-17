package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUrl(t *testing.T) {
	url := NewUrl(
		"url",
	).(*Url)

	assert.Equal(t, "url", url.Url)
	assert.Equal(t, "", url.Host)
	assert.Equal(t, "", url.Port)
	assert.Equal(t, "", url.Path)
	assert.Equal(t, "", url.SearchPart)
}

func TestUrl_Parse(t *testing.T) {
	t.Run("parse url", func(t *testing.T) {
		// Arrange
		expectedUrl := Url{
			Url:        "http://example.com/path?searchPart",
			Host:       "example.com",
			Port:       "80",
			Path:       "path",
			SearchPart: "searchPart",
		}

		// Act
		url, err := NewUrl("http://example.com/path?searchPart").Parse()

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, expectedUrl, url)
	})

	t.Run("if not http schema, return error", func(t *testing.T) {
		// Act
		_, err := NewUrl("https://example.com").Parse()

		// Assert
		assert.ErrorIs(t, ErrOnlyHttpSchemaSupported, err)
	})
}

func TestUrl_isHttp(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "return true, if url begins with http://",
			url:      "http://example.com",
			expected: true,
		},
		{
			name:     "return true, if url contains with http://",
			url:      "hello;http://example.com",
			expected: true,
		},
		{
			name:     "return false, if not contains http://",
			url:      "https://example.com",
			expected: false,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Act & Assert
			assert.Equal(t, it.expected, NewUrl(it.url).(*Url).isHttp())
		})
	}
}

func TestUrl_extractHost(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "extract host from url",
			url:      "http://example.com",
			expected: "example.com",
		},
		{
			name:     "extract host from url w/ port",
			url:      "http://example.com:8080",
			expected: "example.com",
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Act & Assert
			assert.Equal(t, it.expected, NewUrl(it.url).(*Url).extractHost())
		})
	}
}

func TestUrl_extractPort(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "if url doesn't have a port, return 80",
			url:      "http://example.com",
			expected: "80",
		},
		{
			name:     "extract port from url",
			url:      "http://example.com:8080",
			expected: "8080",
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Act & Assert
			assert.Equal(t, it.expected, NewUrl(it.url).(*Url).extractPort())
		})
	}
}

func TestUrl_extractPath(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "if url doesn't have a path, return empty string",
			url:      "http://example.com",
			expected: "",
		},
		{
			name:     "extract path from url",
			url:      "http://example.com/path",
			expected: "path",
		},
		{
			name:     "if url has searchPart, trim it",
			url:      "http://example.com/path?searchPart",
			expected: "path",
		},
		{
			name:     "multiple path case, return full path",
			url:      "http://example.com/path1/path2",
			expected: "path1/path2",
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Act & Assert
			assert.Equal(t, it.expected, NewUrl(it.url).(*Url).extractPath())
		})
	}
}

func TestUrl_extractSearchPart(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "if url doesn't have a path param",
			url:      "http://example.com",
			expected: "",
		},
		{
			name:     "if url doesn't have a searchPart, return empty string",
			url:      "http://example.com/path",
			expected: "",
		},
		{
			name:     "extract searchPart from url",
			url:      "http://example.com/path?searchPart",
			expected: "searchPart",
		},
		{
			name:     "if searchPart has multiple query params, return full text",
			url:      "http://example.com/path?searchPart1&searchPart2",
			expected: "searchPart1&searchPart2",
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Act & Assert
			assert.Equal(t, it.expected, NewUrl(it.url).(*Url).extractSearchPart())
		})
	}
}
