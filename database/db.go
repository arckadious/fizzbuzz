// This package handle connection from the API to database
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/arckadious/fizzbuzz/config"

	"github.com/sirupsen/logrus"
)

const (
	timeout             = 5
	queryDefaultTimeout = 300 //ms
)

// DB class
type DB struct {
	dbConnector *sql.DB
	cf          *config.Config
}

// New constructor DB
func New(cf *config.Config) *DB {
	database := DB{
		cf: cf,
	}

	var err error
	database.dbConnector, err = database.connect()
	if err != nil {
		logrus.Fatal(err)
	}

	return &database

}

func (db *DB) GetDefaultContext() (context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), queryDefaultTimeout*time.Millisecond) //Best practice to set query timeouts for all queries, to avoid slow queries.
	return ctx, cancelFunc
}

// GetConnector returns sql connector
func (db *DB) GetConnector() *sql.DB {
	return db.dbConnector
}

// Shutdown closes the database and prevents new queries from starting. Close then waits for all queries that have started processing on the server to finish.
func (db *DB) Shutdown() {
	if db.dbConnector != nil {
		err := db.dbConnector.Close()
		if err != nil {
			logrus.Warn(err)
			return
		}
		logrus.Info("Shutdown mysql connections OK")
	}
}

// connect establishes connection between API and database
func (db *DB) connect() (dbConnector *sql.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&timeout=%ds",
		db.cf.Database.Username,
		db.cf.Database.Password,
		db.cf.Database.Host,
		db.cf.Database.Port,
		db.cf.Database.Name,
		db.cf.Database.Charset,
		timeout)

	dbConnector, err = sql.Open("mysql", dsn) //The returned DB is safe for concurrent use by multiple goroutines and maintains its own pool of idle connections -> https://pkg.go.dev/database/sql#Open
	if err != nil {
		return
	}

	err = dbConnector.Ping()
	if err != nil {
		return
	}

	dbConnector.SetMaxOpenConns(db.cf.Database.MaxOpenConns)
	dbConnector.SetMaxIdleConns(db.cf.Database.MaxIdleConns)
	dbConnector.SetConnMaxLifetime(db.cf.Database.MaxConnLifeTime * time.Second)

	return
}
