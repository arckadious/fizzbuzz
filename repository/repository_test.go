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
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {

	assert := assert.New(t)
	require := require.New(t)
	cf, err := config.New("../tests/mock/parametersOK.json", *validator.New())
	require.NoError(err)
	db := database.New(cf)
	repo := New(db)
	require.NotEqual(nil, repo)

	//////////////////////////
	// Repository.LogToDB() //
	//////////////////////////

	// empty parameter corID
	if err = repo.LogToDB("", "", "", "", "", ""); assert.Error(err) {
		assert.Equal("Logger coudn't audit data: corID empty", err.Error())
	}

	// logger Type unkwown
	if err = repo.LogToDB("", "", "", "corID", "", ""); assert.Error(err) {
		assert.Equal("Logger coudn't audit data: Logger Type unkwown.", err.Error())
	}

	// send simple request and response
	require.NoError(repo.LogToDB("REQUEST", "", "", "test", "", ""))
	assert.NoError(repo.LogToDB("RESPONSE", "", "", "test", "", ""))

	// check rows, validate host = os.hostname, APP_NAME=fizzbuzz
	var hostname string
	var appname string
	hostnameExpected, _ := os.Hostname()

	assert.NoError(repo.db.GetConnector().QueryRow(`
	SELECT APP_NAME, HOST
	FROM MESSAGES_REQUEST as mr
	LEFT JOIN MESSAGES_RESPONSE as mp ON mr.COR_ID = mp.COR_ID
	WHERE mr.COR_ID = ?`, "test").Scan(&appname, &hostname))

	assert.Equal(cst.AppName, appname)
	assert.Equal(hostnameExpected, hostname)

	// delete test rows
	_, err = repo.db.GetConnector().Exec(`
		DELETE a.*, b.* FROM MESSAGES_REQUEST a 
		LEFT JOIN MESSAGES_RESPONSE b ON b.COR_ID = a.COR_ID
		WHERE a.COR_ID = ?`, "test")
	assert.NoError(err)

	// SQL execution error
	db.Shutdown()

	if err = repo.LogToDB("REQUEST", "", "", "test", "", ""); assert.Error(err) {
		assert.Equal("Logger coudn't send REQUEST data: sql: database is closed", err.Error())
	}

}
