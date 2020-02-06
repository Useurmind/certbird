package main

import (
	// "net/http"
)

func main() {
	caCertConfig := CertConfig{
		FilePath: "ca.pem",
		ValidYears: 10,
		IsCA: true,
	}

	ensureCertificate(&caCertConfig)
}