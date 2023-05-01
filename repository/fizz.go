package repository

import (
	"database/sql"
)

type FizzRepository struct {
	*Repository
}

func NewFizz(repo *Repository) *FizzRepository {
	return &FizzRepository{
		repo,
	}

}
func (rf *FizzRepository) GetMostRequestUsed() (msg string, hits int, noRows bool, err error) {
	err = rf.db.QueryRow("SELECT MSG, count(*) as HITS FROM `MESSAGES_REQUEST` WHERE CHECKSUM IS NOT NULL AND CHECKSUM != '' GROUP BY CHECKSUM ORDER BY HITS DESC LIMIT 1;").Scan(&msg, &hits)

	if err == sql.ErrNoRows {
		noRows = true
	}

	return
}
