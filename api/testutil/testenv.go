package testutil

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/BrunoQuaresma/openticket/api"
)

type TestEnv struct {
	Debug      bool
	HTTPServer *http.Server
	Database   *TestDatabase
	Port       int
	URL        string
}

func (s *TestEnv) Start() {
	logger := io.Discard
	if s.Debug {
		logger = os.Stdout
	}

	dbPort, err := getFreePort()
	if err != nil {
		panic("error getting free port: " + err.Error())
	}
	s.Database = NewTestDatabase(dbPort)
	err = s.Database.Start(logger)
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
	s.URL = "http://localhost:" + fmt.Sprint(s.Port)
}

func (s *TestEnv) Close() {
	s.Database.Stop()
	s.HTTPServer.Close()
}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
