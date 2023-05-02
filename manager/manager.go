// This package process data from action class handlers
package manager

import (
	"github.com/arckadious/fizzbuzz/config"
	"github.com/arckadious/fizzbuzz/repository"
	"github.com/arckadious/fizzbuzz/response"

	"github.com/go-playground/validator/v10"
)

// Manager Class
type Manager struct {
	cf          *config.Config
	apiResponse response.ApiResponse
	validator   *validator.Validate
	repo        *repository.Repository
}

// New constructor Manager
func New(cf *config.Config, apiResponse response.ApiResponse, v *validator.Validate, repo *repository.Repository) *Manager {
	return &Manager{cf, apiResponse, v, repo}
}

// GetApiResponse returns the Api Response
func (m *Manager) GetApiResponse() response.ApiResponse {
	return m.apiResponse
}

// GetValidator returns the validator
func (m *Manager) GetValidator() *validator.Validate {
	return m.validator
}

// GetConfig returns a copy of the config
func (m *Manager) GetConfig() config.Config {
	return *m.cf
}
