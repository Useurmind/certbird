package main

import (
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/useurmind/certbird/caclient"
)

func main() {
	signAddress := "http://localhost:8091/sign"

	defer func() {
		if r:= recover(); r != nil {
			fmt.Println("Panic:", r)
		}
	}()

	csrInfo := caclient.CertRequestInfo{
		DNSNames: []string {"localhost"},
		ValidHours: 1,
	}

	certPackage, err := caclient.CreateCertRequest(csrInfo)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Sending CSR sign request to", signAddress)
	resp, err := http.Post(signAddress, "", bytes.NewBuffer(certPackage.CsrPEM))
	if err != nil{
		panic(err.Error())
	}

	fmt.Println("Response status code is ", resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Status %d - %s: %s", resp.StatusCode, resp.Status, string(body)))
	}

	fmt.Println("Returned certificate:")
	fmt.Println(string(body))
}