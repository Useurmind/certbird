package caserver

import (
	"os"
	"testing"
	"time"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"github.com/useurmind/certbird/utils"
)

func TestEnsureCACertificate(t* testing.T) {
	certFilePath := "./ca_test_cert.pem"
	keyFilePath := "./ca_test_cert.key"
	defer os.Remove(certFilePath)
	defer os.Remove(keyFilePath)

	serverConfig := ServerConfig {
		CACertFilePath: certFilePath,
		CAKeyFilePath: keyFilePath,
	}

	certConfig := &CertConfig{
		IsCA: true,
	}
	certConfig.ParseIPAddresses([]string { "121.12.34.12", "233.23.54.87" })

	err := EnsureCACertificate(certConfig, serverConfig)
	assert.Nil(t, err)

	assert.Equal(t, true, utils.DoesFileExist(certFilePath))
	assert.Equal(t, true, utils.DoesFileExist(keyFilePath))
}

func TestSignedCertIsValid(t *testing.T) {
	dnsName := "myfance.host.com"

	certFilePath := "./ca_test_cert.pem"
	keyFilePath := "./ca_test_cert.key"
	defer os.Remove(certFilePath)
	defer os.Remove(keyFilePath)

	serverConfig := ServerConfig {
		CACertFilePath: certFilePath,
		CAKeyFilePath: keyFilePath,
	}

	validDuration, _ := time.ParseDuration("1h")

	certConfig := &CertConfig{
		IsCA: true,
		ValidDuration: validDuration,
	}
	// certConfig.parseIPAddresses([]string { "121.12.34.12", "233.23.54.87" })

	err := EnsureCACertificate(certConfig, serverConfig)
	assert.Nil(t, err)

	csrPEM := createTestCertificateRequestPEM(dnsName)

	signedCertPEM, err := signCSR(csrPEM, validDuration, certFilePath, keyFilePath)

	caCertPEM, _ := ioutil.ReadFile(certFilePath)
	valid := testCertValid(dnsName, signedCertPEM, caCertPEM)

	assert.Equal(t, true, valid)
}

func createTestCertificateRequestPEM(dnsName string) []byte {

	privKey, _ := rsa.GenerateKey(rand.Reader, 4096)

	// ipAddresses, _ := utils.ConvertIPsFromStringToIP([]string { dnsName })

	certReqInfo := x509.CertificateRequest{
		// IPAddresses: ipAddresses,
		DNSNames: []string { dnsName },
		EmailAddresses: nil,
		PublicKey: privKey.PublicKey,
	}

	csr, _ := x509.CreateCertificateRequest(rand.Reader, &certReqInfo, privKey)

	csrPEM := new(bytes.Buffer)
	pem.Encode(csrPEM, &pem.Block{
		Type:  utils.PEM_TYPE_CERTIFICATE_REQUEST,
		Bytes: csr,
	})

	return csrPEM.Bytes()
}

func testCertValid(dnsName string, signedCertPEM []byte, caCertPEM []byte) bool {
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM([]byte(caCertPEM))
	if !ok {
		panic("failed to parse root certificate")
	}

	block, _ := pem.Decode([]byte(signedCertPEM))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	opts := x509.VerifyOptions{
		Roots: certPool,
		DNSName:       dnsName,
		Intermediates: x509.NewCertPool(),
	}

	if _, err := cert.Verify(opts); err != nil {
		return false
	}

	return true
}