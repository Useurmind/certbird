package caserver

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func GetCACertificate(serverConfig ServerConfig, w http.ResponseWriter, req *http.Request) {
	content, err := ioutil.ReadFile(serverConfig.CACertFilePath)
	if(err != nil) {
		fmt.Println("Could not find CA certificate in", serverConfig.CACertFilePath)
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprintf(w, string(content))
}