package database

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	return testDB.pg.Start()
}

func (testDB *LocalDatabase) Migrate() error {
	appPath := os.Getenv("APP_PATH")
	m, err := migrate.New("file://"+path.Join(appPath, "api", "database", "migrations"), testDB.URL())
	if err != nil {
		testDB.Stop()
		return errors.New("error creating migration instance: " + err.Error())
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
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

func NewLocalDatabase(port uint32, runtimePath string, logger io.Writer) *LocalDatabase {
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

	return testDB
}
