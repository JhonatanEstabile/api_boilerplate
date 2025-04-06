package db

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func GetDBConnection() *sqlx.DB {
	db, err := sqlx.Connect("mysql", "root:root@/api_boilerplate")

	if err != nil {
		log.Fatalln("Error opening database: ", err)
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}
