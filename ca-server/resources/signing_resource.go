package caserver

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func postSignCSR(caFilePath string, w http.ResponseWriter, req *http.Request) {
	caPEM, err := ioutil.ReadFile(caFilePath)
	if(err != nil) {
		fmt.Println("Could not find ca certificate in", caFilePath)
		w.WriteHeader(http.StatusNotFound)
	}

	caPrivKeyPEM, err := ioutil.ReadFile(caFilePath)
	if(err != nil) {
		fmt.Println("Could not find ca certificate in", caFilePath)
		w.WriteHeader(http.StatusNotFound)
	}

	sign
}