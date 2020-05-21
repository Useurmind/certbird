package caserver

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/useurmind/certbird/utils"
)

// Service implements the core functions of the CA server.
type Service struct {
	serverConfig *ServerConfig
}

// NewService creates a new ca server service.
func NewService(serverConfig *ServerConfig) *Service {
	return &Service{
		serverConfig: serverConfig,
	}
}

// GetCACertificate returns the CA certificate (chain) in PEM encoding.
func (s *Service) GetCACertificate() (string, *utils.HTTPError) {
	content, err := ioutil.ReadFile(s.serverConfig.CACertFilePath)
	if(err != nil) {
		log.Println("Could not find CA certificate in", s.serverConfig.CACertFilePath)
		return "", utils.NewHTTPErrorf(http.StatusNotFound, "Could not find CA certificate")
	}
	
	return string(content), nil
}

// SignCSR returns a PEM encoded certificate for the given certificate signing request (csr).
// The parameter csrPEM should contain a PEM encoded CSR that needs to be signed.
// Optionally specify a valid duration that can be parsed with https://golang.org/pkg/time/#ParseDuration.
// If not valid duration is specified a default value will be used.
func (s *Service) SignCSR(csrPEM string, validDuration string, ) (string, *utils.HTTPError) {
	if validDuration == "" {
		validDuration = "30m"
	}

	validDur, err := time.ParseDuration(validDuration)
	if err != nil {
		log.Println("Could not parse valid duration from csr signing request:", err)
		return "", utils.NewHTTPErrorf(http.StatusBadRequest, "Could not parse valid duration %s", validDuration)
	}

	signedCertPEM, err := signCSR([]byte(csrPEM), validDur, s.serverConfig.CACertFilePath, s.serverConfig.CAKeyFilePath)
	if err != nil {
		log.Println("Could not sign certificate request:", err)
		return "", utils.NewHTTPErrorf(http.StatusInternalServerError, "Could not sign certificate request")
	}

	return string(signedCertPEM), nil
}