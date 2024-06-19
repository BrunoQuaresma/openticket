package testutil

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
)

type TestDatabase struct {
	username    string
	password    string
	database    string
	port        uint32
	conn        *embeddedpostgres.EmbeddedPostgres
	runtimePath string
	logger      io.Writer
}

func (testDB *TestDatabase) Start() error {
	err := testDB.conn.Start()
	if err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	go func() {
		<-stop
		err := testDB.Stop()
		if err != nil {
			panic("error shutting down database: " + err.Error())
		} else {
			log.Println("database gracefully stopped")
		}
	}()

	migrateCmd := exec.Command("./scripts/migrate.sh")
	migrateCmd.Env = append(migrateCmd.Env, "POSTGRES_DB_URL="+testDB.URL())
	migrateCmd.Stdout = testDB.logger
	migrateCmd.Stderr = testDB.logger
	err = migrateCmd.Run()
	if err != nil {
		testDB.Stop()
		return errors.New("error running migrations: " + err.Error())
	}

	return nil
}

func (testDB *TestDatabase) Stop() error {
	return testDB.conn.Stop()
}

func (testDB *TestDatabase) URL() string {
	return "postgresql://" + testDB.username + ":" + testDB.password + "@localhost:" + fmt.Sprint(testDB.port) + "/" + testDB.database + "?sslmode=disable"
}

type NewTestDatabaseConfig struct {
	port        int
	runtimePath string
	logger      io.Writer
}

func NewTestDatabase(config NewTestDatabaseConfig) (*TestDatabase, error) {
	testDB := &TestDatabase{
		username:    "postgres",
		password:    "postgres",
		database:    "postgres",
		port:        uint32(config.port),
		runtimePath: config.runtimePath,
	}

	testDB.conn = embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Port(testDB.port).
			Username(testDB.username).
			Password(testDB.password).
			Database(testDB.database).
			Logger(testDB.logger).
			RuntimePath(testDB.runtimePath),
	)

	return testDB, nil
}
