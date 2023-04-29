package config

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	//config ok
	_, err := New("../tests/mock/parametersOK.json", *validator.New())
	assert.Equal(t, nil, err)

	//config file not exist
	_, err = New("../parameters/totooo.json", *validator.New())
	assert.NotEqual(t, nil, err)

	//config file wrong format
	_, err = New("../main.go", *validator.New())
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "invalid character 'p' looking for beginning of value", err.Error())

	//Gin debug mode
	_, err = New("../tests/mock/parameters-gindebug.json", *validator.New())
	assert.Equal(t, nil, err)
	assert.Equal(t, "debug", gin.Mode())
	err = os.Remove("gin.log")
	assert.Equal(t, nil, err)

	//config validation failed
	_, err = New("../tests/mock/parameters-wrongvalue.json", *validator.New())
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "Key: 'Config.Env' Error:Field validation for 'Env' failed on the 'oneof' tag", err.Error())

	//RootPath
	c := &Config{}
	c.InitRootPath("test")
	assert.Equal(t, "test", c.GetRootPath())
	c.InitRootPath("toto")
	assert.Equal(t, "test", c.GetRootPath())

}
