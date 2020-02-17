package main

import (
	"fmt"
	"time"
	"net/http"
	"github.com/useurmind/certbird/caserver"
)

func main() {
	validDuration, _ :=time.ParseDuration("10y")
	caCertConfig := caserver.CertConfig{
		ValidDuration: validDuration,
		IsCA: true,
	}

	serverConfig := caserver.DefaultServerConfig()

	caserver.EnsureCACertificate(&caCertConfig, serverConfig)

	runCertbirdEndpoint(serverConfig)
}

func runCertbirdEndpoint(serverConfig caserver.ServerConfig) {
	http.HandleFunc("/ca", func(w http.ResponseWriter, req *http.Request) { caserver.GetCACertificate(serverConfig, w, req) })
	http.HandleFunc("/sign", func(w http.ResponseWriter, req *http.Request) { caserver.PostSignCSR(serverConfig, w, req) })

	listenOn := ":8091"
	fmt.Println("Listening on", listenOn)
	http.ListenAndServe(listenOn, nil)
}



