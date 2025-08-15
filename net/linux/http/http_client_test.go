package http

import (
	"net"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewHttpClient(t *testing.T) {
	// Act & Assert
	assert.NotNil(t, NewHttpClient())
}

func TestHttpClient_GET(t *testing.T) {
	// Arrange
	client := NewHttpClient()

	/*
		t.Run("if failed to receive response, return ErrFailedToReceiveResponse", func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Mock
			patches.ApplyFunc(
				net.LookupIP,
				func(host string) ([]net.IP, error) {
					return []net.IP{net.IPv4(127, 0, 0, 1)}, nil
				},
			)
			patches.ApplyFunc(
				net.ResolveTCPAddr,
				func(network, address string) (*net.TCPAddr, error) {
					return &net.TCPAddr{
						IP:   net.IPv4(127, 0, 0, 1),
						Port: 80,
						Zone: "",
					}, nil
				},
			)
			patches.ApplyFunc(
				net.DialTCP,
				func(network string, laddr, raddr *net.TCPAddr) (*net.TCPConn, error) {
					return &net.TCPConn{}, nil
				},
			)
			patches.ApplyMethod(
				reflect.TypeOf(&net.TCPConn{}),
				"Write",
				func(_ *net.TCPConn, _ []byte) (int, error) {
					return 10, nil
				},
			)
			patches.ApplyMethod(
				reflect.TypeOf(&net.TCPConn{}),
				"Read",
				func(_ *net.TCPConn, _ []byte) (int, error) {
					return 0, net.UnknownNetworkError("test")
				},
			)

			// Act
			_, err := client.GET("test.example.com", 80, "/users")

			// Assert
			assert.ErrorIs(t, ErrFailedToReceiveResponse, err)
		})
	*/

	t.Run("if failed to send request, return ErrFailedToSendRequest", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Mock
		patches.ApplyFunc(
			net.LookupIP,
			func(host string) ([]net.IP, error) {
				return []net.IP{net.IPv4(127, 0, 0, 1)}, nil
			},
		)
		patches.ApplyFunc(
			net.ResolveTCPAddr,
			func(network, address string) (*net.TCPAddr, error) {
				return &net.TCPAddr{
					IP:   net.IPv4(127, 0, 0, 1),
					Port: 80,
					Zone: "",
				}, nil
			},
		)
		patches.ApplyFunc(
			net.DialTCP,
			func(network string, laddr, raddr *net.TCPAddr) (*net.TCPConn, error) {
				return &net.TCPConn{}, nil
			},
		)
		patches.ApplyMethod(
			reflect.TypeOf(&net.TCPConn{}),
			"Write",
			func(_ *net.TCPConn, _ []byte) (int, error) {
				return 0, net.UnknownNetworkError("test")
			},
		)

		// Act
		_, err := client.GET("test.example.com", 80, "/users")

		// Assert
		assert.ErrorIs(t, ErrFailedToSendRequest, err)
	})

	t.Run("if failed to create connection, return ErrFailedToCreateSocketConnection", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Mock
		patches.ApplyFunc(
			net.LookupIP,
			func(host string) ([]net.IP, error) {
				return []net.IP{net.IPv4(127, 0, 0, 1)}, nil
			},
		)
		patches.ApplyFunc(
			net.ResolveTCPAddr,
			func(network, address string) (*net.TCPAddr, error) {
				return &net.TCPAddr{
					IP:   net.IPv4(127, 0, 0, 1),
					Port: 80,
					Zone: "",
				}, nil
			},
		)
		patches.ApplyFunc(
			net.DialTCP,
			func(network string, laddr, raddr *net.TCPAddr) (*net.TCPConn, error) {
				return nil, net.UnknownNetworkError(network)
			},
		)

		// Act
		_, err := client.GET("test.example.com", 80, "/users")

		// Assert
		assert.ErrorIs(t, ErrFailedToCreateSocketConnection, err)
	})

	t.Run("if failed to create socket address, return ErrFailedToCreateSocketAddress", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Mock
		patches.ApplyFunc(
			net.LookupIP,
			func(host string) ([]net.IP, error) {
				return []net.IP{net.IPv4(127, 0, 0, 1)}, nil
			},
		)
		patches.ApplyFunc(
			net.ResolveTCPAddr,
			func(network, address string) (*net.TCPAddr, error) {
				return nil, net.UnknownNetworkError(network)
			},
		)

		// Act
		_, err := client.GET("test.example.com", 80, "/users")

		// Assert
		assert.ErrorIs(t, ErrFailedToCreateSocketAddress, err)
	})

	t.Run("if 0 length ip address find, return ErrFailedToFindIPAddresses", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Mock
		patches.ApplyFunc(
			net.LookupIP,
			func(host string) ([]net.IP, error) {
				return []net.IP{}, nil
			},
		)

		// Act
		_, err := client.GET("test.example.com", 80, "/users")

		// Assert
		assert.ErrorIs(t, ErrFailedToFindIPAddresses, err)

	})

	t.Run("if failed to look up IP addresses, return ErrFailedToFindIPAddresses", func(t *testing.T) {
		// Act
		_, err := client.GET("test.example.com", 80, "/users")

		// Assert
		assert.ErrorIs(t, ErrFailedToFindIPAddresses, err)
	})
}
