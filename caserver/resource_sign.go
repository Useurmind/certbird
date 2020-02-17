package caserver

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"time"
)

func PostSignCSR(serverConfig ServerConfig, w http.ResponseWriter, req *http.Request) {
	csrPEM, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Could not read body for csr signing request:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var validDurationString string
	validDurationStrings, ok := req.URL.Query()["validDuration"]
	if !ok || len(validDurationStrings[0]) < 1 {
		validDurationString = "30m"
	} else {
		validDurationString = validDurationStrings[0]
	}

	validDuration, err := time.ParseDuration(validDurationString)
	if err != nil {
		log.Println("Could not parse valid duration from csr signing request:", err)
		fmt.Fprintln(w, "Invalid duration string:", validDurationString)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	signedCertPEM, err := signCSR(csrPEM, validDuration, serverConfig.CACertFilePath, serverConfig.CAKeyFilePath)
	if err != nil {
		log.Println("Could not sign certificate request:", err)
		fmt.Fprintln(w, "Invalid certificate request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(signedCertPEM)
}