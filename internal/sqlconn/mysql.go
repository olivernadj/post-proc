package sqlconn

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func GetConnection()(*sql.DB, error) {
	return sql.Open("mysql", os.Getenv("DATA_SOURCE_NAME"))
}
