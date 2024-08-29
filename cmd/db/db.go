package db

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLStorage(config mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}

	// initialize database
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("DB: Connected Successfully")
	return db, nil
}
