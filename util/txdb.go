package util

import (
	"database/sql"
	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func init() {
	txdb.Register("txdb", "mysql", "root:root@tcp(localhost:3306)/storage")
}

func InitDb() (*sql.DB, error) {
	db, err := sql.Open("txdb", uuid.New().String())

	if err != nil {
		return db, err
	}

	return db, db.Ping()
}
