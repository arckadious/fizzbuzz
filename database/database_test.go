// This package connect the API to database
package database

import (
	"testing"

	"github.com/arckadious/fizzbuzz/config"
	"github.com/arckadious/fizzbuzz/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	cf, err := config.New("../tests/mock/parametersOK.json", *validator.New())
	assert.Equal(t, nil, err)
	Connect(cf)
}
