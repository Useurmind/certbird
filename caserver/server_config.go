package caserver

// ServerConfig defines the configuration properties for the ca server.
type ServerConfig struct {
	CACertFilePath string
	CAKeyFilePath string
}

// DefaultServerConfig returns the default server configuration.
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		CACertFilePath: "ca.pem",
		CAKeyFilePath: "ca.key",
	}
}