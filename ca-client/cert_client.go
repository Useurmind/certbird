package caclient

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"crypto/rsa"
	"encoding/pem"

	"github.com/useurmind/certbird/utils"
)

// CertificatePackage holds all information a ca client/web server needs about its certificate.
type CertificatePackage struct {
	PrivKey *rsa.PrivateKey
	CertPEM []byte
	CsrPEM []byte
	CertBundlePEM []byte
}

// CertRequestInfo can be used to specify the target of an CRS.
type CertRequestInfo struct {
	DNSNames []string
	IPAddresses []string
	EmailAddresses []string
	
	ValidYears int32
	ValidMonths int32
	ValidDays int32
	ValidHours int32
}

// Returns the csr, cert, and priv key (PEM encoded) and potential error.
func CreateCertRequest(csrInfo CertRequestInfo) (*CertificatePackage, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	ipAddresses, err := utils.ConvertIPsFromStringToIP(csrInfo.IPAddresses)
	if err != nil {
		return nil, err
	}

	certReqInfo := x509.CertificateRequest{
		IPAddresses: ipAddresses,
		DNSNames: csrInfo.DNSNames,
		EmailAddresses: csrInfo.EmailAddresses,
		PublicKey: privKey.PublicKey,
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, &certReqInfo, privKey)
	if err != nil {
		return nil, err
	}

	csrPEM := new(bytes.Buffer)
	pem.Encode(csrPEM, &pem.Block{
		Type:  utils.PEM_TYPE_CERTIFICATE_REQUEST,
		Bytes: csr,
	})

	return &CertificatePackage{
		PrivKey: privKey,
		CsrPEM: csrPEM.Bytes(),
	}, nil
}