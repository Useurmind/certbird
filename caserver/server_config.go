package caserver

import (
	"fmt"
)

// ServerConfig defines the configuration properties for the ca server.
type ServerConfig struct {
	// The address to listen on, use empty for all
	Address string

	// The port ot listen on
	Port int

	// The path to the PEM encoded CA certificate file.
	CACertFilePath string

	// The path to the PEM encoded CA certificate private key file.
	CAKeyFilePath string
}

// DefaultServerConfig returns the default server configuration.
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		CACertFilePath: "ca.pem",
		CAKeyFilePath:  "ca.key",
	}
}

// TestServerConfig returns the server configuration for testing purposes.
func TestServerConfig(basePath string) ServerConfig {
	return ServerConfig{
		Address:        "127.0.0.1",
		Port:           8091,
		CACertFilePath: fmt.Sprintf("%s/ca.pem", basePath),
		CAKeyFilePath:  fmt.Sprintf("%s/ca.key", basePath),
	}
}
