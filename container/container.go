// This package init all classes
package container

import (
	fizzaction "github.com/arckadious/fizzbuzz/action/fizz"
	"github.com/arckadious/fizzbuzz/repository"

	"net/http"

	cst "github.com/arckadious/fizzbuzz/constant"

	"github.com/arckadious/fizzbuzz/config"
	"github.com/arckadious/fizzbuzz/database"
	"github.com/arckadious/fizzbuzz/manager"
	"github.com/arckadious/fizzbuzz/response"

	"github.com/go-playground/validator/v10"
)

// Container class
type Container struct {
	Conf      *config.Config
	Validator *validator.Validate

	Manager     *manager.Manager
	ManagerFizz *manager.Fizz

	FizzAction *fizzaction.FizzAction

	Db       *database.DB
	Repo     *repository.Repository
	RepoFizz *repository.Fizz
}

// New constructor Container
func New(conf *config.Config, validator *validator.Validate, db *database.DB) *Container {

	container := Container{
		Conf:      conf,
		Validator: validator,
		Db:        db,
	}
	container.setRepositories(db)
	container.setManagers()
	container.setActions()

	return &container
}

func (c *Container) setManagers() {

	// Default API Response is always set to HTTP status 200
	resp := response.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     cst.StatusSuccess,
		Messages:   make([]response.ApiError, 0),
	}

	// Starting init managers and actions
	c.Manager = manager.New(c.Conf, resp, c.Validator, c.Repo) //init manager object by calling constructor

	c.ManagerFizz = manager.NewFizz(c.Manager, c.RepoFizz) //init fizz object, child of manager class
}
func (c *Container) setActions() {
	c.FizzAction = fizzaction.New(c.ManagerFizz)
}

func (c *Container) setRepositories(db *database.DB) {
	c.Repo = repository.New(db)
	c.RepoFizz = repository.NewFizz(c.Repo)
}
