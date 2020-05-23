package caserver

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/useurmind/certbird/caclient"
)

func TestGetCACertReturnsCACert(t *testing.T) {
	ctx := NewTestContext(t)
	serverConfig := ctx.EnsureCertificate()
	defer ctx.CleanupTest()

	service := NewService(serverConfig)

	caCertPEM, err := service.GetCACertificate()
	assert.Nil(t, err)

	assert.NotEqual(t, "", caCertPEM)
}

func TestGetCACertReturns404IfCACertMissing(t *testing.T) {
	ctx := NewTestContext(t)
	serverConfig := ctx.EnsureCertificate()
	defer ctx.CleanupTest()

	os.Remove(serverConfig.CACertFilePath)

	service := NewService(serverConfig)

	_, err := service.GetCACertificate()
	assert.NotNil(t, err)
	assert.Equal(t, 404, err.HTTPStatusCode)
}

func TestSignCSRReturnsValidCert(t *testing.T) {
	ctx := NewTestContext(t)
	serverConfig := ctx.EnsureCertificate()
	defer ctx.CleanupTest()

	dnsName := "myfancy.host.com"
	csrInfo := caclient.CertRequestInfo{
		DNSNames: []string{dnsName},
	}

	csrPkg, err := caclient.CreateCertRequest(csrInfo)
	assert.Nil(t, err)

	service := NewService(serverConfig)

	certPEM, herr := service.SignCSR(string(csrPkg.CsrPEM), "1h")
	assert.Nil(t, herr)

	err = ctx.ValidateCertificatePEM(certPEM, dnsName)
	assert.Nil(t, err)
}
