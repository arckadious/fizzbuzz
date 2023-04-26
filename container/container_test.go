package container

import (
	"testing"

	"github.com/arckadious/fizzbuzz/config"
	"github.com/arckadious/fizzbuzz/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	cf, err := config.New("../parameters/parameters.json", *validator.New())
	assert.Equal(t, nil, err)
	container := New(cf, validator.New())
	assert.NotEqual(t, nil, container)
	assert.NotEqual(t, nil, container.Conf)
	assert.NotEqual(t, nil, container.Validator)
	assert.NotEqual(t, nil, container.FizzAction)
	assert.NotEqual(t, nil, container.Repo)
	assert.NotEqual(t, nil, container.RepoFizz)
}
