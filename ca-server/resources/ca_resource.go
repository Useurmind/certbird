package caserver

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func getCaCertificate(caFilePath string, w http.ResponseWriter, req *http.Request) {
	content, err := ioutil.ReadFile(caFilePath)
	if(err != nil) {
		fmt.Println("Could not find ca certificate in", caFilePath)
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprintf(w, string(content))
}