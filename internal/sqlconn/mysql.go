package sqlconn

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func GetConnection()(*sql.DB, error) {
	return sql.Open("mysql", "root:example@tcp(db:3306)/postproc?parseTime=true")
}
