package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kenza-ai/kenza/logutil"
	_ "github.com/lib/pq" // postgres driver init
	"github.com/pressly/goose"
)

const pathToMigrations = "/kenza/data/db/migrations"

// New - Connects to, pings and returns a sql DB
func New(user, password, host, dbName string, port int64, migrate bool) (*sql.DB, error) {
	var err error
	var db *sql.DB
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbName)

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for range ticker.C {
		if db, err = sql.Open("postgres", connectionString); err == nil {
			break
		}
		logutil.Info("opening db '%s' failed '%s', retrying in 5\"", dbName, err)
	}

	for range ticker.C {
		if err = db.Ping(); err == nil {
			break
		}
		logutil.Info("pinging db '%s' failed '%s', retrying in 5\"", dbName, err)
	}

	if !migrate {
		return db, nil
	}

	if err := goose.Run("up", db, pathToMigrations); err != nil {
		return nil, err
	}

	return db, nil
}
