package container

import (
	fizzaction "github.com/arckadious/fizzbuzz/action/fizz"
	"github.com/arckadious/fizzbuzz/database"
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

	Repo     *repository.Repository
	RepoFizz *repository.RepositoryFizz
}

// Container constructor
func New(conf *config.Config, validator *validator.Validate) *Container {

	container := Container{
		Conf:      conf,
		Validator: validator,
	}
	container.setRepositories()
	container.setActions()

	return &container
}

// Set application actions
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

func (c *Container) setRepositories() {
	c.Repo = repository.NewRepository(database.New(c.Conf))
	c.RepoFizz = repository.NewFizzRepository(c.Repo)
}
