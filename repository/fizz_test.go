// This package contains all functions which interact with database
package repository

import (
	"testing"

	"github.com/arckadious/fizzbuzz/config"
	"github.com/arckadious/fizzbuzz/database"
	"github.com/arckadious/fizzbuzz/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFizz(t *testing.T) {

	assert := assert.New(t)
	require := require.New(t)
	cf, err := config.New("../tests/mock/parametersOK.json", *validator.New())
	require.NoError(err)
	db := database.New(cf)
	repo := New(db)
	require.NotNil(repo)
	repoFizz := NewFizz(repo)
	require.NotNil(repoFizz)

	///////////////////////////////
	// Fizz.GetMostRequestUsed() //
	///////////////////////////////

	// delete test rows
	_, err = repo.db.GetConnector().Exec(`
	DELETE FROM MESSAGES_REQUEST`)
	require.NoError(err)

	// send simple request logs
	assert.NoError(repo.LogToDB("REQUEST", "test", "", "testFizz", "test", ""))
	assert.NoError(repo.LogToDB("REQUEST", "test", "", "testFizz2", "test", ""))
	assert.NoError(repo.LogToDB("REQUEST", "test", "", "testFizz3", "test", ""))
	assert.NoError(repo.LogToDB("REQUEST", "test", "", "testFizz4", "test", ""))
	assert.NoError(repo.LogToDB("REQUEST", "test2", "", "testFizz5", "test2", ""))
	assert.NoError(repo.LogToDB("REQUEST", "test2", "", "testFizz6", "test2", ""))

	msg, hits, noRows, err := repoFizz.GetMostRequestUsed()
	assert.NoError(err)
	assert.False(noRows)
	assert.Equal(4, hits)
	assert.Equal("test", msg)

	// delete test rows
	_, err = repo.db.GetConnector().Exec(`
		DELETE FROM MESSAGES_REQUEST`)
	assert.NoError(err)

	// No rows check
	msg, hits, noRows, err = repoFizz.GetMostRequestUsed()
	assert.Error(err)
	assert.True(noRows)
	assert.Equal(0, hits)
	assert.Equal("", msg)

	// SQL execution error
	db.Shutdown()

	msg, hits, noRows, err = repoFizz.GetMostRequestUsed()
	assert.Error(err)
	assert.False(noRows)
	assert.Equal(0, hits)
	assert.Equal("", msg)

}
