package testutil

import (
	"io"
	"net"
	"os"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/BrunoQuaresma/openticket/sdk"
	"github.com/brianvoe/gofakeit"
)

type Credentials struct {
	Email    string
	Password string
}

type TestEnv struct {
	Debug            bool
	Database         *TestDatabase
	Server           *api.Server
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
	tEnv.Database = NewTestDatabase(dbPort)
	err = tEnv.Database.Start(logger)
	if err != nil {
		panic("error starting test database: " + err.Error())
	}

	port, err := getFreePort()
	if err != nil {
		panic("error getting free port: " + err.Error())
	}
	tEnv.Server = api.Start(api.Options{
		DatabaseURL: tEnv.Database.URL(),
		Mode:        api.TestMode,
		Port:        port,
	})
}

func (tEnv *TestEnv) Close() {
	tEnv.Database.Stop()
	tEnv.Server.Close()
}

func (tEnv *TestEnv) URL() string {
	return "http://localhost" + tEnv.Server.Addr()
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

	_, err := tEnv.SDK().Setup((api.SetupRequest{
		Name:     gofakeit.Name(),
		Username: gofakeit.Username(),
		Email:    credentials.Email,
		Password: credentials.Password,
	}))
	if err != nil {
		panic("error making setup request: " + err.Error())
	}
}

func (tEnv *TestEnv) AdminCredentials() Credentials {
	return tEnv.adminCredentials
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
