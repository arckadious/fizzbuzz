// This package init all classes
package container

import (
	"testing"

	"github.com/arckadious/fizzbuzz/config"
	"github.com/arckadious/fizzbuzz/database"
	"github.com/arckadious/fizzbuzz/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	assert := assert.New(t)

	/////////////////////
	// Container.New() //
	/////////////////////

	cf, err := config.New("../tests/mock/parametersOK.json", *validator.New())
	assert.NoError(err)
	container := New(cf, validator.New(), &database.DB{})
	assert.NotEqual(nil, container)
	assert.NotEqual(nil, container.Conf)
	assert.NotEqual(nil, container.Validator)
	assert.NotEqual(nil, container.FizzAction)
	assert.NotEqual(nil, container.Repo)
	assert.NotEqual(nil, container.RepoFizz)
	assert.Equal(&database.DB{}, container.Db)
}
