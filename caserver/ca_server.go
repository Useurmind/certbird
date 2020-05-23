package caserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// CAServer can be run to setup the CA server.
type CAServer struct {
	ServerConfig ServerConfig

	// The underlying http server that listens for requests.
	server *http.Server

	// Wait group used when running the server async.
	wg *sync.WaitGroup
}

func (s *CAServer) RunAsync() error {
	s.wg = &sync.WaitGroup{}

	service, err := s.init()
	if err != nil {
		return err
	}

	err = s.createServer(service)
	if err != nil {
		return err
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		err := s.runServer()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ERROR: Server crashed when running async: %v", err)
		}

		log.Printf("Server shutdown done")
	}()

	return nil
}

func (s *CAServer) Shutdown() error {
	if s.wg == nil {
		return fmt.Errorf("Cannot use Shutdown when RunAsync was not called before")
	}

	log.Printf("Shutting down CAServer\r\n")
	err := s.server.Shutdown(context.Background())
	if err != nil {
		return err
	}

	s.wg.Wait()

	return nil
}

func (s *CAServer) Run() error {
	service, err := s.init()
	if err != nil {
		return err
	}

	err = s.createServer(service)
	if err != nil {
		return err
	}

	err = s.runServer()
	if err != nil {
		return err
	}

	return nil
}

func (s *CAServer) init() (*Service, error) {
	validDuration, _ := time.ParseDuration("10y")
	caCertConfig := CertConfig{
		ValidDuration: validDuration,
		IsCA: true,
	}

	err := EnsureCACertificate(&caCertConfig, s.ServerConfig)
	if err != nil {
		return nil, err
	}

	service := NewService(&s.ServerConfig)

	return service, nil
}

func (s *CAServer) createServer(service *Service) error {
	if s.server != nil {
		return fmt.Errorf("Server was already created")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/ca", func(w http.ResponseWriter, req *http.Request) { GetCACertificate(service, w, req) })
	mux.HandleFunc("/sign", func(w http.ResponseWriter, req *http.Request) { PostSignCSR(service, w, req) })

	listenOn := fmt.Sprintf("%s:%d", s.ServerConfig.Address, s.ServerConfig.Port)
	s.server = &http.Server{
		Addr: listenOn,
		Handler: mux,
	}

	return nil
}

func (s *CAServer) runServer() error {
	log.Println("Server starting, listening on", s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}