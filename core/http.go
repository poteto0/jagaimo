package core

import (
	"fmt"
	"strconv"
	"strings"
)

type Header struct {
	Name  string
	Value string
}

func NewHeader(name, value string) Header {
	return Header{
		Name:  name,
		Value: value,
	}
}

type HttpResponse struct {
	Version    string
	StatusCode uint32
	Reason     string
	Headers    []Header
	Body       string
}

func NewHttpResponse(rawResponse string) (HttpResponse, error) {
	preprocessedResponse := preprocessedResponseBody(rawResponse)
	statusLine, remaining, ok := splitToStrPair(preprocessedResponse, "\n")
	if !ok {
		return HttpResponse{}, NewUnexpectedInputError(
			fmt.Sprintf("invalid http response {%s}", rawResponse),
			nil,
		)
	}

	rawHeaders, body, ok := splitToStrPair(remaining, "\n\n")
	if !ok {
		body = remaining
	}

	headers := []Header{}
	for rawHeader := range strings.SplitSeq(rawHeaders, "\n") {
		splitHeader := strings.SplitN(rawHeader, ":", 2)
		if len(splitHeader) < 2 {
			return HttpResponse{}, NewUnexpectedInputError(
				fmt.Sprintf("invalid http headers {%s}", rawHeader),
				nil,
			)
		}

		headers = append(
			headers,
			NewHeader(
				strings.TrimLeft(splitHeader[0], " "),
				strings.TrimLeft(splitHeader[1], " "),
			),
		)
	}

	statues := strings.Split(statusLine, " ")
	if len(statues) < 3 {
		return HttpResponse{}, NewUnexpectedInputError(
			fmt.Sprintf("invalid http status line {%s}", statusLine),
			nil,
		)
	}

	statusCode, ok := parseStatusCode(statues[1])
	if !ok {
		statusCode = 404
	}

	return HttpResponse{
		Version:    statues[0],
		StatusCode: statusCode,
		Reason:     strings.Join(statues[2:], " "),
		Headers:    headers,
		Body:       body,
	}, nil
}

func (res *HttpResponse) HeaderValue(name string) (string, error) {
	for _, header := range res.Headers {
		if header.Name == name {
			return header.Value, nil
		}
	}

	return "", NewUnexpectedInputError(
		fmt.Sprintf("header {%s} not found", name),
		nil,
	)
}

func preprocessedResponseBody(rawResponse string) string {
	return strings.ReplaceAll(
		strings.TrimLeft(rawResponse, " "),
		"\r\n", "\n",
	)
}

func splitToStrPair(s, sep string) (string, string, bool) {
	parts := strings.SplitN(s, sep, 2)
	if len(parts) < 2 {
		return s, "", false
	}

	return parts[0], parts[1], true
}

func parseStatusCode(statusCode string) (uint32, bool) {
	code, err := strconv.ParseUint(statusCode, 10, 32)
	if err != nil {
		return 0, false
	}

	return uint32(code), true
}
