package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func main() {
	var caFilePath string = "ca.pem"

	caCertConfig := CertConfig{
		FilePath: caFilePath,
		ValidYears: 10,
		IsCA: true,
	}

	ensureCertificate(&caCertConfig)

	runCaEndpoint(caFilePath)
}

func runCaEndpoint(caFilePath string) {
	http.HandleFunc("/ca", func(w http.ResponseWriter, req *http.Request) { getCaCertificate(caFilePath, w, req) })

	listenOn := ":8091"
	fmt.Println("Listening on", listenOn)
	http.ListenAndServe(listenOn, nil)
}

func getCaCertificate(caFilePath string, w http.ResponseWriter, req *http.Request) {
	content, err := ioutil.ReadFile(caFilePath)
	if(err != nil) {
		fmt.Println("Could not find ca certificate in", caFilePath)
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprintf(w, string(content))
}

