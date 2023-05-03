// This package load API configuration
package config

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {

	assert := assert.New(t)

	//////////////////
	// Config.New() //
	//////////////////

	//config ok
	_, err := New("../tests/mock/parametersOK.json", *validator.New())
	assert.NoError(err)

	//config file not exist
	_, err = New("../parameters/totooo.json", *validator.New())
	assert.Error(err)

	//config file wrong format
	_, err = New("../main.go", *validator.New())
	if assert.Error(err) {
		assert.Equal("invalid character '/' looking for beginning of value", err.Error())
	}

	//config validation failed
	_, err = New("../tests/mock/parameters-wrongvalue.json", *validator.New())
	if assert.Error(err) {
		assert.Equal("Key: 'Config.Env' Error:Field validation for 'Env' failed on the 'oneof' tag", err.Error())
	}

	///////////////////////
	// Config.RootPath() //
	///////////////////////

	c := &Config{}
	c.InitRootPath("test")
	assert.Equal("test", c.GetRootPath())
	c.InitRootPath("toto")
	assert.Equal("test", c.GetRootPath())

}
