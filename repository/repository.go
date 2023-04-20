package repository

import "database/sql"

type Repository struct {
	db *sql.DB
}

// NewRepository
func NewRepository(Db *sql.DB) *Repository {
	return &Repository{Db}
}
