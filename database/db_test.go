// This package connect the API to database
package database

import (
	"testing"

	"github.com/arckadious/fizzbuzz/config"
	"github.com/arckadious/fizzbuzz/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {

	assert := assert.New(t)

	///////////////////////
	// DB.GetConnector() //
	///////////////////////

	cf, err := config.New("../tests/mock/parametersOK.json", *validator.New())
	assert.NoError(err)

	db := New(cf)
	assert.NotEqual(nil, db)
	assert.NotEqual(nil, db.GetConnector()) //check if it's nil value //TO DO
	db.Shutdown()

	//db.connect() // mock test to do

	//benchmark test to do
}
