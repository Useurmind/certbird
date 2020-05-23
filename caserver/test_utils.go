package caserver

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestContext holds common information required for most tests.
type TestContext struct {
	// storageFolder is the folder holding all test data
	// there should be nothing in there except the data for the current test
	// putting all files into a deletable folder that makes cleanup easier
	storageFolder string

	serverConfig ServerConfig

	t *testing.T
}

func NewTestContext(t *testing.T) *TestContext {
	storageFolder := "test_data"
	return &TestContext{
		storageFolder: storageFolder,
		serverConfig: TestServerConfig(storageFolder),
		t:             t,
	}
}

func (c *TestContext) PrepareTest() *ServerConfig {
	validDuration, _ := time.ParseDuration("1h")

	err := os.MkdirAll(c.storageFolder, 666)
	assert.Nil(c.t, err)

	certConfig := &CertConfig{
		IsCA:          true,
		ValidDuration: validDuration,
	}
	err = EnsureCACertificate(certConfig, c.serverConfig)
	assert.Nil(c.t, err)

	return &c.serverConfig
}

func (c *TestContext) CleanupTest() {
	os.RemoveAll(c.storageFolder)
}

// ValidateCertificatePEM checks if the certificate is valid for the given dnsname.
// Returns an error if not.
func (c *TestContext) ValidateCertificatePEM(certPEM string, dnsName string) error {
	caCertPEM, err := ioutil.ReadFile(c.serverConfig.CACertFilePath)
	if err != nil {
		return err
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM([]byte(caCertPEM))
	if !ok {
		panic("failed to parse root certificate")
	}

	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return fmt.Errorf("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %s", err)
	}

	opts := x509.VerifyOptions{
		Roots: certPool,
		DNSName:       dnsName,
		Intermediates: x509.NewCertPool(),
	}

	if _, err := cert.Verify(opts); err != nil {
		return fmt.Errorf("Certificate is not valid")
	}

	return nil
}
