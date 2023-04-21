package repository

import (
	"database/sql"
	"os"
	"strings"

	cst "github.com/arckadious/fizzbuzz/constant"

	"github.com/sirupsen/logrus"
)

type Repository struct {
	db *sql.DB
}

// NewRepository
func NewRepository(Db *sql.DB) *Repository {
	return &Repository{Db}
}

func (r *Repository) LogToDB(logType, msg, url, corID, checksum, status string) {
	if strings.ToUpper(logType) != "REQUEST" && strings.ToUpper(logType) != "RESPONSE" {
		logrus.Error("Logger coudn't audit data: Logger Type unkwown.")
		return
	} else if corID == "" {
		logrus.Error("Logger coudn't audit data: corID empty -", corID)
		return
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
		logrus.Error("Logger coudn't send data: Logger Type (request,response ?) unkwown.")
		return
	}

	_, err := r.db.Exec(sql, vals...)
	if err != nil {
		logrus.Error("Logger coudn't send "+logType+" data: ", err)
		return
	}
}
