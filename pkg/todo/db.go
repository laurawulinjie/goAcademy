package todo

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	connStr := "host=localhost port=5432 user=laura.wu dbname=todoapp sslmode=disable"
	var err error

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return DB.Ping()
}
