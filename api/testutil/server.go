package testutil

import (
	"io"
	"net/http"
	"os"

	"github.com/BrunoQuaresma/openticket/api"
)

type TestServer struct {
	Debug      bool
	HTTPServer *http.Server
	Database   *TestDatabase
	Port       int
}

func (s *TestServer) Start() {
	logger := io.Discard
	if s.Debug {
		logger = os.Stdout
	}

	s.Database = NewTestDatabase()
	err := s.Database.Start(logger)
	if err != nil {
		panic("error starting test database: " + err.Error())
	}

	s.Port, err = getFreePort()
	if err != nil {
		panic("error getting free port: " + err.Error())
	}
	s.HTTPServer = api.Start(api.Options{
		DatabaseURL: s.Database.URL(),
		Debug:       s.Debug,
		Port:        s.Port,
	})
}

func (s *TestServer) Close() {
	s.Database.Stop()
	s.HTTPServer.Close()
}
