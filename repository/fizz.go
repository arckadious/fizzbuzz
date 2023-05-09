// This package contains all functions which interact with database
package repository

import (
	"database/sql"
)

// Fizz class (repository child class)
type Fizz struct {
	*Repository
}

type FizzInterface interface { // Use an interface as prototype, to allow mocks testing.
	GetMostRequestUsed() (msg string, hits int, noRows bool, err error)
}

// New constructor Repository child Fizz
func NewFizz(repo *Repository) *Fizz {
	return &Fizz{
		repo,
	}

}

// GetMostRequestUsed returns from database the most used request (if any) as well as the number of hits.
func (rf *Fizz) GetMostRequestUsed() (msg string, hits int, noRows bool, err error) {
	ctx, cancel := rf.db.GetDefaultContext()
	defer cancel()

	err = rf.db.GetConnector().QueryRowContext(ctx, "SELECT MAX(MSG), count(*) as HITS FROM `MESSAGES_REQUEST` WHERE CHECKSUM IS NOT NULL AND CHECKSUM != '' GROUP BY CHECKSUM ORDER BY HITS DESC LIMIT 1;").Scan(&msg, &hits)
	// NOTE - MAX(MSG) instead of MSG: Since version 5.7.5 Mysql has the "ONLY_FULL_GROUP_BY" flag by default enabled. This request select a random 'MSG' value, but the new flag doesn't allow it. See https://jira.mariadb.org/browse/MDEV-10426 for more details

	if err == sql.ErrNoRows {
		noRows = true
	}

	return
}
