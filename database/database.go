// This package connect the API to database
package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/arckadious/fizzbuzz/config"

	"github.com/sirupsen/logrus"
)

// Connect connect the API to database
func Connect(cf *config.Config) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cf.Database.Username,
		cf.Database.Password,
		cf.Database.Host,
		cf.Database.Port,
		cf.Database.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logrus.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		logrus.Fatal(err)
	}

	db.SetMaxOpenConns(cf.Database.MaxOpenConns)
	db.SetMaxIdleConns(cf.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cf.Database.MaxConnLifeTime * time.Second)

	return db
}
