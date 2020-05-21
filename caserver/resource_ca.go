package caserver

import (
	"fmt"
	"net/http"
)

func GetCACertificate(service *Service, w http.ResponseWriter, req *http.Request) {
	content, err := service.GetCACertificate()
	if err != nil {
		err.WriteTo(w)
	}
	fmt.Fprintf(w, content)
}