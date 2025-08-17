package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHeader(t *testing.T) {
	// Act & Assert
	assert.Equal(t, Header{"name", "value"}, NewHeader("name", "value"))
}

func TestNewHttpResponse(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected HttpResponse
		}{
			{
				name:  "invalid code to 404",
				input: "HTTP/1.1 xxx Not Found\n hello: world\n\nhello world !!",
				expected: HttpResponse{
					Version:    "HTTP/1.1",
					StatusCode: 404,
					Reason:     "Not Found",
					Headers: []Header{
						NewHeader("hello", "world"),
					},
					Body: "hello world !!",
				},
			},
			{
				name:  "valid code to uint32",
				input: "HTTP/1.1 200 OK\n hello: world\n\nhello world !!",
				expected: HttpResponse{
					Version:    "HTTP/1.1",
					StatusCode: 200,
					Reason:     "OK",
					Headers: []Header{
						NewHeader("hello", "world"),
					},
					Body: "hello world !!",
				},
			},
		}

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				// Act
				response, err := NewHttpResponse(it.input)

				// Assert
				assert.NoError(t, err)
				assert.Equal(t, it.expected, response)
			})
		}
	})

	t.Run("error case", func(t *testing.T) {
		tests := []struct {
			name    string
			input   string
			message string
		}{
			{
				name:    "invalid http response",
				input:   "HTTP/1.1",
				message: "invalid http response {HTTP/1.1}",
			},
			{
				name:    "invalid http headers",
				input:   "HTTP/1.1 \ninvalid header",
				message: "invalid http headers {invalid header}",
			},
			{
				name:    "invalid http stats line",
				input:   "HTTP/1.1 \n hello: world",
				message: "invalid http status line {HTTP/1.1 }",
			},
		}

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				// Act & Assert
				_, err := NewHttpResponse(it.input)
				assert.Error(t, err)
				assert.Equal(t, it.message, err.Error())
			})
		}
	})
}

func TestHttpResponse_HeaderValue(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		// Arrange
		response := HttpResponse{
			Headers: []Header{
				NewHeader("hello", "world"),
			},
		}

		// Act
		val, err := response.HeaderValue("hello")

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, "world", val)
	})

	t.Run("error case", func(t *testing.T) {
		// Arrange
		response := HttpResponse{
			Headers: []Header{
				NewHeader("hello", "world"),
			},
		}

		// Act & Assert
		_, err := response.HeaderValue("invalid")
		assert.Error(t, err)
		assert.Equal(t, "header {invalid} not found", err.Error())
	})
}

func Test_preprocessedResponseBody(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "trim left spaces",
			input:    "   hello world !! ",
			expected: "hello world !! ",
		},
		{
			name:     "replace \r\n with \n",
			input:    "hello\r\nworld",
			expected: "hello\nworld",
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Act & Assert
			assert.Equal(t, it.expected, preprocessedResponseBody(it.input))
		})
	}
}

func Test_splitToStrPair(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		sep            string
		expectedFirst  string
		expectedSecond string
		expectedOk     bool
	}{
		{
			name:           "success separated",
			input:          "hello\nworld",
			sep:            "\n",
			expectedFirst:  "hello",
			expectedSecond: "world",
			expectedOk:     true,
		},
		{
			name:           "success separate just once",
			input:          "hello\nworld\n!!",
			sep:            "\n",
			expectedFirst:  "hello",
			expectedSecond: "world\n!!",
			expectedOk:     true,
		},
		{
			name:           "fail separated",
			input:          "hello",
			sep:            "\n",
			expectedFirst:  "hello",
			expectedSecond: "",
			expectedOk:     false,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Act
			first, second, ok := splitToStrPair(it.input, it.sep)

			// Assert
			assert.Equal(t, it.expectedFirst, first)
			assert.Equal(t, it.expectedSecond, second)
			assert.Equal(t, it.expectedOk, ok)
		})
	}
}

func Test_parseStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   uint32
		expectedOk bool
	}{
		{
			name:       "success",
			input:      "200",
			expected:   200,
			expectedOk: true,
		},
		{
			name:       "fail",
			input:      "hello",
			expected:   0,
			expectedOk: false,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Act
			code, ok := parseStatusCode(it.input)

			// Assert
			assert.Equal(t, it.expected, code)
			assert.Equal(t, it.expectedOk, ok)
		})
	}
}
