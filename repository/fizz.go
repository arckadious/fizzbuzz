package repository

import (
	"database/sql"
)

type RepositoryFizz struct {
	*Repository
}

func NewFizzRepository(repo *Repository) *RepositoryFizz {
	return &RepositoryFizz{
		repo,
	}

}
func (rf *RepositoryFizz) GetMostRequestUsed() (msg string, hits int, noRows bool, err error) {
	err = rf.db.QueryRow("SELECT MSG, count(*) as HITS FROM `MESSAGES_REQUEST` WHERE CHECKSUM IS NOT NULL AND CHECKSUM != '' GROUP BY CHECKSUM ORDER BY HITS DESC LIMIT 1;").Scan(&msg, &hits)

	if err == sql.ErrNoRows {
		noRows = true
	}

	return
}
