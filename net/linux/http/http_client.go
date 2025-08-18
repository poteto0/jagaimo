package http

import (
	"errors"
	"net"
	"strconv"

	"github.com/poteto0/jagaimo/core"
)

const MaxResponseSize = 4096

var (
	ErrFailedToFindIPAddresses        = errors.New("failed to find IP addresses")
	ErrFailedToCreateSocketAddress    = errors.New("failed to create socket address")
	ErrFailedToCreateSocketConnection = errors.New("failed to create socket connection")
	ErrFailedToSendRequest            = errors.New("failed to send request")
	ErrFailedToReceiveResponse        = errors.New("failed to receive response")
)

type IHttpClient interface {
	/*
		[host] --> [ip] --> [socketAddr]
	*/
	GET(host string, port uint16, path string) (core.HttpResponse, error)
}

type HttpClient struct{}

func NewHttpClient() IHttpClient {
	return &HttpClient{}
}

func (c *HttpClient) GET(host string, port uint16, path string) (core.HttpResponse, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return core.HttpResponse{}, ErrFailedToFindIPAddresses
	}

	if len(ips) < 1 {
		return core.HttpResponse{}, ErrFailedToFindIPAddresses
	}

	addr, err := net.ResolveTCPAddr("tcp", "["+ips[0].String()+"]"+":"+strconv.Itoa(int(port)))
	if err != nil {
		return core.HttpResponse{}, ErrFailedToCreateSocketAddress
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return core.HttpResponse{}, ErrFailedToCreateSocketConnection
	}
	defer conn.Close()

	// create request line to send -> connection
	request := "GET " + path + " HTTP/1.1\n"
	request += "Host: " + host + "\n"
	request += "Accept: text/html\n"
	request += "Connection: close\n\n"

	if _, err := conn.Write([]byte(request)); err != nil {
		return core.HttpResponse{}, ErrFailedToSendRequest
	}

	// read response
	received := make([]byte, MaxResponseSize)
	if _, err = conn.Read(received); err != nil {
		return core.HttpResponse{}, ErrFailedToReceiveResponse
	}

	// create response
	res, err := core.NewHttpResponse(string(received))
	if err != nil {
		return core.HttpResponse{}, err
	}

	return res, nil
}
