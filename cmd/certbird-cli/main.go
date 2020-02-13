package main

import (
	"github.com/useurmind/certbird/caclient"
)

func main() {
	csrInfo := core.CertRequestInfo{
		DNSNames: []string {"localhost"},
		ValidHours: 1,
	}

	certPackage, err := core.CreateCertRequest(csrInfo)
	if err != nil {
		panic(err.Error())
	}

	println(string(certPackage.CsrPEM))
}