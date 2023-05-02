// This package contains all functions which interact with database
package repository

import (
	"errors"
	"os"
	"strings"

	cst "github.com/arckadious/fizzbuzz/constant"
	"github.com/arckadious/fizzbuzz/database"
)

// Repository class
type Repository struct {
	db *database.DB
}

// New constructor Repository
func New(Db *database.DB) *Repository {
	return &Repository{Db}
}

// LogToDB stores requests and responses in database.
func (r *Repository) LogToDB(logType, msg, url, corID, checksum, status string) (err error) {
	if corID == "" {
		return errors.New("Logger coudn't audit data: corID empty")
	}

	hostname, _ := os.Hostname()
	sql := ""
	vals := []interface{}{}

	switch strings.ToUpper(logType) {
	case "REQUEST":
		sql = "INSERT INTO `MESSAGES_REQUEST` (`MSG`, `COR_ID`, `HOST`, `APP_NAME`, `SERVICE_ADDRESS`, `CHECKSUM`) VALUES (?, ?, ?, ?, ?, ?);"
		vals = append(vals, msg, corID, hostname, cst.AppName, url, checksum)
		break
	case "RESPONSE":
		sql = "INSERT INTO `MESSAGES_RESPONSE` (`MSG`, `COR_ID`, `STATUS`) VALUES (?, ?, ?);"
		vals = append(vals, msg, corID, status)
		break
	default:
		return errors.New("Logger coudn't audit data: Logger Type unkwown.")
	}

	_, err = r.db.GetConnector().Exec(sql, vals...)
	if err != nil {
		return errors.New("Logger coudn't send " + logType + " data: " + err.Error())
	}
	return
}
