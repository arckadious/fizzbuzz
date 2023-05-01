// This package init all classes
package container

import (
	"database/sql"

	fizzaction "github.com/arckadious/fizzbuzz/action/fizz"
	"github.com/arckadious/fizzbuzz/repository"

	"net/http"

	"github.com/arckadious/fizzbuzz/config"
	"github.com/arckadious/fizzbuzz/manager"
	"github.com/arckadious/fizzbuzz/response"

	"github.com/go-playground/validator/v10"
)

type Container struct {
	Conf      *config.Config
	Validator *validator.Validate

	FizzAction *fizzaction.FizzAction
	db         *sql.DB
	Repo       *repository.Repository
	RepoFizz   *repository.RepositoryFizz
}

// New constructor Container
func New(conf *config.Config, validator *validator.Validate, db *sql.DB) *Container {

	container := Container{
		Conf:      conf,
		Validator: validator,
		db:        db,
	}
	container.setRepositories(db)
	container.setActions()

	return &container
}

func (c *Container) setActions() *Container {

	// Default API Response is always set to HTTP status 200
	resp := response.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     response.StatusSuccess,
		Messages:   make([]response.ApiError, 0),
	}

	// Starting init managers and actions
	mng := manager.New(c.Conf, resp, c.Validator, c.Repo) //init manager object by calling constructor

	mngrFizz := manager.NewFizz(mng, c.RepoFizz) //init fizz object, child of manager class

	c.FizzAction = fizzaction.NewFizzAction(mngrFizz)

	return c
}

func (c *Container) setRepositories(db *sql.DB) {
	c.Repo = repository.NewRepository(db)
	c.RepoFizz = repository.NewFizzRepository(c.Repo)
}
