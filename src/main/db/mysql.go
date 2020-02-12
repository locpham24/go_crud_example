package db

import (
	"database/sql"
	"main/conf"
)

var mysql *sql.DB

func GetInstance() *sql.DB {
	if mysql == nil {
		dbUser := conf.USER
		dbPass := conf.PASSWORD
		dbName := conf.DB

		db, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName)
		if err != nil {
			panic(err.Error())
		}
		mysql = db
	}
	return mysql
}
