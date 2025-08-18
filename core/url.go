package core

import (
	"errors"
	"strings"
)

const httpSchema = "http://"
const httpPort = "80"

var (
	ErrOnlyHttpSchemaSupported = errors.New("only HTTP schema is supported")
)

type IUrl interface {
	/*
		Url Analysis
			- only HTTP schema supported
	*/
	Parse() (Url, error)

	isHttp() bool

	extractHost() string

	/*
		if url doesn't have a port,
		return httpPort(80)
	*/
	extractPort() string

	extractPath() string

	extractSearchPart() string
}

type Url struct {
	Url  string
	Host string
	Port string
	Path string

	// query parameter
	SearchPart string
}

func NewUrl(url string) IUrl {
	return &Url{
		Url: url,
	}
}

func (url *Url) Parse() (Url, error) {
	if !url.isHttp() {
		return Url{}, ErrOnlyHttpSchemaSupported
	}

	return Url{
		Url:        url.Url,
		Host:       url.extractHost(),
		Port:       url.extractPort(),
		Path:       url.extractPath(),
		SearchPart: url.extractSearchPart(),
	}, nil
}

func (url *Url) isHttp() bool {
	return strings.Contains(url.Url, httpSchema)
}

func (url *Url) extractHost() string {
	urlParts := extractUrlParts(url.Url, httpSchema)

	// if host has port, trim it
	if portIndex := strings.Index(urlParts[0], ":"); portIndex != -1 {
		return urlParts[0][:portIndex]
	}

	return urlParts[0]
}

func (url *Url) extractPort() string {
	urlParts := extractUrlParts(url.Url, httpSchema)

	// if host has port, trim it
	if portIndex := strings.Index(urlParts[0], ":"); portIndex != -1 {
		return urlParts[0][portIndex+1:]
	}
	return httpPort
}

func (url *Url) extractPath() string {
	urlParts := extractUrlParts(url.Url, httpSchema)

	if len(urlParts) < 2 {
		return ""
	}

	pathAndSearchPart := strings.SplitN(urlParts[1], "?", 2)
	return pathAndSearchPart[0]
}

func (url *Url) extractSearchPart() string {
	urlParts := extractUrlParts(url.Url, httpSchema)

	if len(urlParts) < 2 {
		return ""
	}

	pathAndSearchPart := strings.SplitN(urlParts[1], "?", 2)
	if len(pathAndSearchPart) < 2 {
		return ""
	}

	return pathAndSearchPart[1]
}

//nolint:unparam
func extractUrlParts(url, schema string) []string {
	index := strings.LastIndex(url, schema)
	return strings.SplitN(url[index+len(schema):], "/", 2)
}
