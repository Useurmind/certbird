package caserver

// ServerConfig defines the configuration properties for the ca server.
type ServerConfig struct {
	// The path to the PEM encoded CA certificate file.
	CACertFilePath string

	// The path to the PEM encoded CA certificate private key file.
	CAKeyFilePath string
}

// DefaultServerConfig returns the default server configuration.
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		CACertFilePath: "ca.pem",
		CAKeyFilePath: "ca.key",
	}
}