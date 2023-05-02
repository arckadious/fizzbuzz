// This package contains all functions which interact with database
package repository

import (
	"database/sql"
)

// Fizz class (repository child class)
type Fizz struct {
	*Repository
}

// New constructor Repository child Fizz
func NewFizz(repo *Repository) *Fizz {
	return &Fizz{
		repo,
	}

}

// GetMostRequestUsed returns from database the most used request (if any) as well as the number of hits.
func (rf *Fizz) GetMostRequestUsed() (msg string, hits int, noRows bool, err error) {
	err = rf.db.GetConnector().QueryRow("SELECT MSG, count(*) as HITS FROM `MESSAGES_REQUEST` WHERE CHECKSUM IS NOT NULL AND CHECKSUM != '' GROUP BY CHECKSUM ORDER BY HITS DESC LIMIT 1;").Scan(&msg, &hits)

	if err == sql.ErrNoRows {
		noRows = true
	}

	return
}
