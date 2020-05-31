package main

import (
	"time"
	
	certbirdapi "github.com/useurmind/certbird/k8s-controller/pkg/apis/cr/v1"
	_ "github.com/useurmind/certbird/k8s-controller/pkg/client/clientset/versioned"
	certbirdclientfake "github.com/useurmind/certbird/k8s-controller/pkg/client/clientset/versioned/fake"
	certbirdinformers "github.com/useurmind/certbird/k8s-controller/pkg/client/informers/externalversions"
)

func main() {
	certRequest := certbirdapi.CertRequest{}
	testClient := certbirdclientfake.NewSimpleClientset(&certRequest)
	resyncTime, _ := time.ParseDuration("5s")
	certbirdinformers.NewSharedInformerFactory(testClient, resyncTime)

}