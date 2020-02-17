package caserver

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"time"
	"encoding/json"

	"github.com/useurmind/certbird/utils"
)

// CertConfig configures how a certificate should be created
type CertConfig struct {
	Organization  string
	Country       string
	Province      string
	Locality      string
	StreetAddress string
	PostalCode    string

	DNSNames       []string
	IPAddresses    []net.IP
	EmailAddresses []string

	IsCA          bool
	ValidDuration time.Duration
}

func (c *CertConfig) ParseIPAddresses(ipAddresses []string) error {
	ipAddresses2, err := utils.ConvertIPsFromStringToIP(ipAddresses)
	if err != nil {
		return err
	}

	c.IPAddresses = ipAddresses2

	return nil
}

// EnsureCACertificate makes sure that CA certificate and private key are available.
// If not they will be created from the settings in the serverConfig.
func EnsureCACertificate(config *CertConfig, serverConfig ServerConfig) error {
	log.Println("Ensuring presence of CA certificate and private key")
	// ensure and load the CA private key
	spk := StoredPrivateKey{
		filePath: serverConfig.CAKeyFilePath,
	}
	spk.ensure()

	if !utils.DoesFileExist(serverConfig.CACertFilePath) {
		log.Println("Creating missing CA certificate file", serverConfig.CACertFilePath)
		cert, err := createCertificate(config, spk.privKey, &spk.privKey.PublicKey, nil)
		if err != nil {
			return err
		}

		ioutil.WriteFile(serverConfig.CACertFilePath, cert, 0777)
	}

	return nil
}

// createCertificate creates a certificate for the given config
// parentCa will be used to sign the new certificate
// returns the pem encoded cert and private key bytes
func createCertificate(config *CertConfig, privKey interface{}, publicKey interface{}, parentCa *x509.Certificate) ([]byte, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{config.Organization},
			Country:       []string{config.Country},
			Province:      []string{config.Province},
			Locality:      []string{config.Locality},
			StreetAddress: []string{config.StreetAddress},
			PostalCode:    []string{config.PostalCode},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(config.ValidDuration),
		IsCA:                  config.IsCA,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		DNSNames:              config.DNSNames,
		IPAddresses:           config.IPAddresses,
		EmailAddresses:        config.EmailAddresses,
	}

	if parentCa == nil {
		parentCa = ca
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, parentCa, publicKey, privKey)
	if err != nil {
		return nil, err
	}

	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  utils.PEM_TYPE_CERTIFICATE,
		Bytes: caBytes,
	})
	if err != nil {
		return nil, err
	}

	return caPEM.Bytes(), nil
}

func decodeCertificateRequestFromPEM(csrBytes []byte) (*x509.CertificateRequest, error) {
	csrBlock, _ := pem.Decode(csrBytes)
	if csrBlock == nil {
		return nil, fmt.Errorf("Could not read csr block from PEM")
	}

	if csrBlock.Type != utils.PEM_TYPE_CERTIFICATE_REQUEST {
		return nil, fmt.Errorf("The PEM does not contain a certificate request, found type %s", csrBlock.Type)
	}

	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return csr, nil
}

func loadCertificateFromPEM(filePath string) (*x509.Certificate, error) {
	certBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	certBlock, _ := pem.Decode(certBytes)
	if certBlock == nil {
		return nil, fmt.Errorf("Could not read cert block from PEM file %s", filePath)
	}

	if certBlock.Type != utils.PEM_TYPE_CERTIFICATE {
		return nil, fmt.Errorf("The PEM file %s does not contain a certificate, found type %s", filePath, certBlock.Type)
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func signCSR(csrPEM []byte, validDuration time.Duration, caCertFilePath string, caKeyFilePath string) ([]byte, error) {
	log.Println("Signing certificate request with duration", validDuration.String())
	log.Println("Certificate request (PEM) is:")
	log.Println(string(csrPEM))

	spk := StoredPrivateKey{
		filePath: caKeyFilePath,
	}
	err := spk.loadFromPEM()
	if err != nil {
		return nil, err
	}

	caCert, err := loadCertificateFromPEM(caCertFilePath)
	if err != nil {
		return nil, err
	}

	csr, err := decodeCertificateRequestFromPEM(csrPEM)
	if err != nil {
		return nil, err
	}

	err = csr.CheckSignature()
	if err != nil {
		return nil, err
	}

	certConfig := &CertConfig{
		IsCA: false,

		DNSNames:       csr.DNSNames,
		IPAddresses:    csr.IPAddresses,
		EmailAddresses: csr.EmailAddresses,

		ValidDuration: validDuration,
	}

	certConfigJSON, err := json.Marshal(certConfig)
	if err != nil {
		return nil, err
	}
	log.Println("Certificate request content is:", string(certConfigJSON))

	signedCertPEM, err := createCertificate(certConfig, spk.privKey, csr.PublicKey, caCert)
	if err != nil {
		return nil, err
	}

	return signedCertPEM, nil
}
