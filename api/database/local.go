package database

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

type LocalDatabase struct {
	username string
	password string
	database string
	port     uint32
	pg       *embeddedpostgres.EmbeddedPostgres
	logger   io.Writer
}

func (testDB *LocalDatabase) Start() error {
	err := testDB.pg.Start()
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

func (testDB *LocalDatabase) Stop() error {
	return testDB.pg.Stop()
}

func (testDB *LocalDatabase) URL() string {
	return "postgresql://" + testDB.username + ":" + testDB.password + "@localhost:" + fmt.Sprint(testDB.port) + "/" + testDB.database + "?sslmode=disable"
}

func NewLocalDatabase(port uint32, runtimePath string, logger io.Writer) (*LocalDatabase, error) {
	testDB := &LocalDatabase{
		username: "postgres",
		password: "postgres",
		database: "postgres",
		port:     port,
		logger:   logger,
	}

	testDB.pg = embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Port(testDB.port).
			Username(testDB.username).
			Password(testDB.password).
			Database(testDB.database).
			Logger(testDB.logger).
			RuntimePath(runtimePath),
	)

	return testDB, nil
}
