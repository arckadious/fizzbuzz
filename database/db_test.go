// This package connect the API to database
package database

import (
	"testing"

	"github.com/arckadious/fizzbuzz/config"
	"github.com/arckadious/fizzbuzz/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {

	assert := assert.New(t)
	cf, err := config.New("../tests/mock/parametersOK.json", *validator.New())
	assert.NoError(err)
	hook := new(test.Hook)
	logrus.AddHook(hook)
	logrus.SetLevel(logrus.DebugLevel)

	//////////////
	// DB.New() //
	//////////////
	db := New(cf)

	///////////////////////
	// DB.GetConnector() //
	///////////////////////

	// getConnector OK
	assert.NotNil(db)
	assert.NotNil(db.GetConnector())

	///////////////////
	// DB.shutdown() //
	///////////////////

	// getConnector database closed
	db.Shutdown()
	if assert.NotNil(hook.LastEntry()) {
		assert.Equal("Shutdown mysql connections OK", hook.LastEntry().Message)
	}

	//////////////////
	// DB.connect() //
	//////////////////

	//bad username
	cf.Database.Username = "toto"
	if _, err := db.connect(); assert.Error(err) {
		assert.Contains(err.Error(), "Access denied for user ")
	}

	////////////////////////////
	// DB.GetDefaultContext() //
	////////////////////////////
	ctx, cancel := db.GetDefaultContext()
	assert.NotNil(ctx)
	assert.NotNil(cancel)
	defer cancel()

	_, ok := ctx.Deadline()
	assert.True(ok)

}
