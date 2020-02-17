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

	serverConfig := caserver.DefaultServerConfig()

	EnsureCACertificate(&caCertConfig, serverConfig)

	runCertbirdEndpoint(serverConfig)
}

func runCertbirdEndpoint(serverConfig caserver.ServerConfig) {
	http.HandleFunc("/ca", func(w http.ResponseWriter, req *http.Request) { caserver.GetCACertificate(serverConfig, w, req) })
	http.HandleFunc("/sign", func(w http.ResponseWriter, req *http.Request) { caserver.PostSignCSR(serverConfig, w, req) })

	listenOn := ":8091"
	fmt.Println("Listening on", listenOn)
	http.ListenAndServe(listenOn, nil)
}



