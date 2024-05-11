package testutils

import (
	"io"
	"os"

	server "github.com/BrunoQuaresma/openticket"
)

type TestServerConfig struct {
	Debug bool
}

type Stop func()

func RunTestServer(c TestServerConfig) Stop {
	logger := io.Discard
	if c.Debug {
		logger = os.Stdout
	}

	testDB := NewTestDatabase()
	err := testDB.Start(logger)
	if err != nil {
		panic("error starting test database: " + err.Error())
	}

	s := server.Start(server.Options{
		DatabaseURL: testDB.URL(),
		Debug:       c.Debug,
	})

	return func() {
		testDB.Stop()
		s.Close()
	}
}
