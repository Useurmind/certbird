package caserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func PostSignCSR(service *Service, w http.ResponseWriter, req *http.Request) {
	csrPEM, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintln(w, "Could not read csr from body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validDuration := ""
	validDurationStrings, ok := req.URL.Query()["validDuration"]
	if ok && len(validDurationStrings[0]) > 0 {
		validDuration = validDurationStrings[0]
	}
	
	signedCertPEM, e := service.SignCSR(string(csrPEM), validDuration)
	if e != nil {
		e.WriteTo(w)
		return
	}

	w.Write([]byte(signedCertPEM))
}