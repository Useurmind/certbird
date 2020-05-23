package main

import (
	"log"
	"os"

	"github.com/useurmind/certbird/caserver"
)

func main() {
	server := caserver.CAServer{
		ServerConfig: caserver.DefaultServerConfig(),
	}

	err := server.Run()	
	if err != nil {
		log.Printf("ERROR while running caserver: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
