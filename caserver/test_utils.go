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

	ServerConfig ServerConfig

	server *CAServer

	t *testing.T
}

func NewTestContext(t *testing.T) *TestContext {
	storageFolder := "test_data"
	return &TestContext{
		storageFolder: storageFolder,
		ServerConfig:  TestServerConfig(storageFolder),
		t:             t,
	}
}

func (c *TestContext) EnsureCertificate() *ServerConfig {
	validDuration, _ := time.ParseDuration("1h")

	certConfig := &CertConfig{
		IsCA:          true,
		ValidDuration: validDuration,
	}
	err := EnsureCACertificate(certConfig, c.ServerConfig)
	assert.Nil(c.t, err)

	return &c.ServerConfig
}

func (c *TestContext) StartServer() {
	c.server = &CAServer{
		ServerConfig: c.ServerConfig,
	}

	err := c.server.RunAsync()
	assert.Nil(c.t, err)
}

func (c *TestContext) CleanupTest() {
	if c.server != nil {
		err := c.server.Shutdown()
		assert.Nil(c.t, err)
	}

	os.RemoveAll(c.storageFolder)
}

// ValidateCertificatePEM checks if the certificate is valid for the given dnsname.
// Returns an error if not.
func (c *TestContext) ValidateCertificatePEM(certPEM string, dnsName string) error {
	caCertPEM, err := ioutil.ReadFile(c.ServerConfig.CACertFilePath)
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
		Roots:         certPool,
		DNSName:       dnsName,
		Intermediates: x509.NewCertPool(),
	}

	if _, err := cert.Verify(opts); err != nil {
		return fmt.Errorf("Certificate is not valid")
	}

	return nil
}
