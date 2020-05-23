package caclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// CAServerClient can be used to request certificates from the CA server.
type CAServerClient struct {
	// CAServers is a list of CA server instances.
	// You can put in valid urls to reach the CA server instances.
	// e.g. https://mycaserver.local:8080
	CAServers []string

	// The DNS names that should be requested in the certificate.
	RequestedDNSNames []string

	// The duration the requested SSL certificate should be valid.
	ValidDuration time.Duration
}

// RequestCert can be used to request an SSL certificate from any of the CAServers
// that are configured.
func (c *CAServerClient) RequestCert() (string, error) {
	csrInfo := CertRequestInfo{
		DNSNames: c.RequestedDNSNames,
	}
	csrPkg, err := CreateCertRequest(csrInfo)
	if err != nil {
		return "", err
	}

	certPem := ""
	for _, caserverAddr := range c.CAServers {
		url := fmt.Sprintf("%s/sign?validDuration=%s", caserverAddr, c.ValidDuration.String())
		resp, err := http.Post(url, "application/x-pem-file", strings.NewReader(string(csrPkg.CsrPEM)))
		if err != nil {
			continue
		}

		certBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		certPem = string(certBytes)

		break
	}

	if certPem == "" {
		return "", fmt.Errorf("Could retrieve certificate from caserver, last error: %v", err)
	}

	return certPem, nil
}

