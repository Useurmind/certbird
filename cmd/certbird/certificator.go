package main

import (
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
)

// CertConfig configures how a certificate should be created
type CertConfig struct {
	Organization string
	Country string
	Province string
	Locality string
	StreetAddress string
	PostalCode string

	
	FilePath string
	IsCA bool
	ValidYears int32
	ValidMonths int32
	ValidDays int32
	ValidHours int32
}


func ensureCertificate(config *CertConfig) error {
	if _, err := os.Stat(config.FilePath); os.IsNotExist(err) {
		cert, privKey, err := createCertificate(config)
		if err != nil {
			return err
		}

		ioutil.WriteFile(config.FilePath, cert, 0777)
		certExt := filepath.Ext(config.FilePath)
		privKeyPath := config.FilePath[0:len(config.FilePath) - len(certExt)] + ".key"
		ioutil.WriteFile(privKeyPath, privKey, 0777)
	}

	return nil
}

// returns the pem encoded cert and private key bytes
func createCertificate(config *CertConfig) ([]byte, []byte, error)  {
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
		NotAfter:              time.Now().AddDate(int(config.ValidYears), int(config.ValidMonths), int(config.ValidDays)),
		IsCA:                  config.IsCA,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
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