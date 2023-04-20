package manager

import (
	"github.com/arckadious/fizzbuzz/config"
	"github.com/arckadious/fizzbuzz/repository"
	"github.com/arckadious/fizzbuzz/response"

	"github.com/go-playground/validator/v10"
)

// Manager
type Manager struct {
	cf          *config.Config
	apiResponse response.ApiResponse
	validator   *validator.Validate
	repo        *repository.Repository
}

// Constructor NewManager
func New(cf *config.Config, apiResponse response.ApiResponse, v *validator.Validate, repo *repository.Repository) *Manager {
	return &Manager{cf, apiResponse, v, repo}
}

// Returns the ApiResponse
func (m *Manager) GetApiResponse() response.ApiResponse {
	return m.apiResponse
}

func (m *Manager) GetValidator() *validator.Validate {
	return m.validator
}

// Returns a copy of the config
func (m *Manager) GetConfig() config.Config {
	return *m.cf
}
