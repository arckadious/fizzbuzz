// This package contains all functions which interact with database
package repository

import (
	"database/sql"
)

type Fizz struct {
	*Repository
}

func NewFizz(repo *Repository) *Fizz {
	return &Fizz{
		repo,
	}

}
func (rf *Fizz) GetMostRequestUsed() (msg string, hits int, noRows bool, err error) {
	err = rf.db.GetConnector().QueryRow("SELECT MSG, count(*) as HITS FROM `MESSAGES_REQUEST` WHERE CHECKSUM IS NOT NULL AND CHECKSUM != '' GROUP BY CHECKSUM ORDER BY HITS DESC LIMIT 1;").Scan(&msg, &hits)

	if err == sql.ErrNoRows {
		noRows = true
	}

	return
}
