// This package contains all functions which interact with database
package repository

import (
	"os"
	"testing"

	"github.com/arckadious/fizzbuzz/config"
	cst "github.com/arckadious/fizzbuzz/constant"
	"github.com/arckadious/fizzbuzz/database"
	"github.com/arckadious/fizzbuzz/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {

	cf, err := config.New("../tests/mock/parametersOK.json", *validator.New())
	assert.Equal(t, nil, err)
	db := database.New(cf)
	repo := New(db)
	assert.NotEqual(t, nil, repo)

	// empty parameter corID
	err = repo.LogToDB("", "", "", "", "", "")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "Logger coudn't audit data: corID empty", err.Error())

	// logger Type unkwown
	err = repo.LogToDB("", "", "", "corID", "", "")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "Logger coudn't audit data: Logger Type unkwown.", err.Error())

	// send simple request and response
	err = repo.LogToDB("REQUEST", "", "", "test", "", "")
	assert.Equal(t, nil, err)
	err = repo.LogToDB("RESPONSE", "", "", "test", "", "")
	assert.Equal(t, nil, err)

	// check rows, validate host = os.hostname, APP_NAME=fizzbuzz
	var hostname string
	var appname string
	hostnameExpected, _ := os.Hostname()

	err = repo.db.GetConnector().QueryRow(`
		SELECT APP_NAME, HOST
		FROM MESSAGES_REQUEST as mr
		LEFT JOIN MESSAGES_RESPONSE as mp ON mr.COR_ID = mp.COR_ID
		WHERE mr.COR_ID = ?
	`, "test").Scan(&appname, &hostname)
	assert.Equal(t, nil, err)
	assert.Equal(t, cst.AppName, appname)
	assert.Equal(t, hostnameExpected, hostname)

	// delete test rows
	_, err = repo.db.GetConnector().Exec(`
		DELETE a.*, b.* FROM MESSAGES_REQUEST a 
		LEFT JOIN MESSAGES_RESPONSE b ON b.COR_ID = a.COR_ID
		WHERE a.COR_ID = ?`, "test")
	assert.Equal(t, nil, err)

	// SQL execution error
	db.Shutdown()
	err = repo.LogToDB("REQUEST", "", "", "test", "", "")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "Logger coudn't send REQUEST data: sql: database is closed", err.Error())

}
