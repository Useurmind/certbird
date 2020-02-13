package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/useurmind/certbird/caserver"
)

func main() {
	var caFilePath string = "ca.pem"

	caCertConfig := CertConfig{
		ValidYears: 10,
		IsCA: true,
	}

	ensureCertificate(&caCertConfig, caFilePath)

	runCaEndpoint(caFilePath)
}

func runCaEndpoint(caFilePath string) {
	http.HandleFunc("/ca", func(w http.ResponseWriter, req *http.Request) { caserver.getCaCertificate(caFilePath, w, req) })

	listenOn := ":8091"
	fmt.Println("Listening on", listenOn)
	http.ListenAndServe(listenOn, nil)
}



