package testutil

import (
	"io"
	"net"
	"os"

	"github.com/BrunoQuaresma/openticket/api"
	database "github.com/BrunoQuaresma/openticket/api/database/gen"
	"github.com/BrunoQuaresma/openticket/sdk"
	"github.com/brianvoe/gofakeit"
)

type Credentials struct {
	Email    string
	Password string
}

type TestEnv struct {
	Debug            bool
	database         *TestDatabase
	server           *api.Server
	sdk              *sdk.Client
	adminCredentials Credentials
}

func (tEnv *TestEnv) Start() {
	logger := io.Discard
	if tEnv.Debug {
		logger = os.Stdout
	}

	dbPort, err := getFreePort()
	if err != nil {
		panic("error getting free port: " + err.Error())
	}
	tEnv.database = NewTestDatabase(dbPort)
	err = tEnv.database.Start(logger)
	if err != nil {
		panic("error starting test database: " + err.Error())
	}

	port, err := getFreePort()
	if err != nil {
		panic("error getting free port: " + err.Error())
	}
	tEnv.server = api.Start(api.Options{
		DatabaseURL: tEnv.database.URL(),
		Mode:        api.TestMode,
		Port:        port,
	})
}

func (tEnv *TestEnv) Close() {
	tEnv.database.Stop()
	tEnv.server.Close()
}

func (tEnv *TestEnv) URL() string {
	return "http://localhost" + tEnv.server.Addr()
}

func (tEnv *TestEnv) SDK() *sdk.Client {
	if tEnv.sdk == nil {
		tEnv.sdk = sdk.New(tEnv.URL())
	}
	return tEnv.sdk
}

func (tEnv *TestEnv) Setup() {
	credentials := Credentials{
		Email:    gofakeit.Email(),
		Password: FakePassword(),
	}
	tEnv.adminCredentials = credentials

	var res api.Response[any]
	_, err := tEnv.SDK().Setup(api.SetupRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    credentials.Email,
		Password: credentials.Password,
	}, &res)
	if err != nil {
		panic("error making setup request: " + err.Error())
	}
}

func (tEnv *TestEnv) AdminCredentials() Credentials {
	return tEnv.adminCredentials
}

func (tEnv *TestEnv) DBQueries() *database.Queries {
	return tEnv.server.DBQueries()
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
