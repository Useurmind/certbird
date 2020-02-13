package caserver

import (
	"fmt"
	"net"
	"bytes"
	"os"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"crypto/rsa"
	"math/big"
	"time"
	"io/ioutil"
	"encoding/pem"
	"path/filepath"
	. "github.com/ahmetb/go-linq"
)

// CertConfig configures how a certificate should be created
type CertConfig struct {
	Organization string
	Country string
	Province string
	Locality string
	StreetAddress string
	PostalCode string

	DNSNames []string
	IPAddresses []string
	EmailAddresses []string
	
	IsCA bool
	ValidYears int32
	ValidMonths int32
	ValidDays int32
	ValidHours int32
}


func ensureCertificate(config *CertConfig, certFilePath string) error {
	if _, err := os.Stat(certFilePath); os.IsNotExist(err) {
		cert, privKey, err := createCertificate(config, nil)
		if err != nil {
			return err
		}

		ioutil.WriteFile(certFilePath, cert, 0777)
		certExt := filepath.Ext(certFilePath)
		privKeyPath := certFilePath[0:len(certFilePath) - len(certExt)] + ".key"
		ioutil.WriteFile(privKeyPath, privKey, 0777)
	}

	return nil
}

// returns the pem encoded cert and private key bytes
func createCertificate(config *CertConfig, privKey *rsa.PrivateKey) ([]byte, []byte, error)  {
	ipAddresses := convertIPsFromStringToIP(config.IPAddresses)

	durationHours, err := time.ParseDuration(fmt.Sprintf("%dh", config.ValidHours))
	if err != nil {
		return nil, nil, err
	}

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
		NotAfter:              time.Now().AddDate(int(config.ValidYears), int(config.ValidMonths), int(config.ValidDays)).Add(durationHours),
		IsCA:                  config.IsCA,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		DNSNames: config.DNSNames,
		IPAddresses: ipAddresses,
		EmailAddresses: config.EmailAddresses,
	}

	caPrivKey := privKey

	if caPrivKey == nil {
		var err error
		caPrivKey, err = rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, nil, err
		}
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	
	return caPEM.Bytes(), caPrivKeyPEM.Bytes(), nil
}

func signCSR(csrPEM []byte, caPEM []byte, caPrivKeyPEM[])
{
	
}

func convertIPsFromStringToIP(ipAddressStrings []string) ([]net.IP) {
	var ipAddresses []net.IP
	From(ipAddressStrings).Select(func(ipString interface{}) interface{} { return net.ParseIP(ipString.(string)) }).ToSlice(&ipAddresses)

	return ipAddresses
}