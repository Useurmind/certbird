package main

import (
	"net/http"

	"github.com/useurmind/certbird/caclient"
)

func main() {
	csrInfo := caclient.CertRequestInfo{
		DNSNames: []string {"localhost"},
		ValidHours: 1,
	}

	certPackage, err := caclient.CreateCertRequest(csrInfo)
	if err != nil {
		panic(err.Error())
	}

	resp, err := http.Post("localhost:8091/sign", "", certPackage.CsrPEM)
	if err != nil{
		panic(err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Status %d - %s: %s", resp.StatusCode, resp.Status, resp.Body)
	}

	println(string(resp.Body))
}