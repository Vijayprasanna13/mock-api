package databases

import (
	"database/sql"
	"log"
	"mock-api/config"
	_ "mock-api/go-sql-driver/mysql"
)

func OpenConn() *sql.DB {

	/*
	*----------------------------------------------------
	*					DATABASE CREDENTIALS
	*----------------------------------------------------
	*Set the Database credentials
	 */

	DB_DATABASE := config.DB_DATABASE
	DB_USERNAME := config.DB_USERNAME
	DB_PASSWORD := config.DB_PASSWORD
	DB_HOST := config.DB_HOST
	DB_PORT := config.DB_PORT

	/*
	*-----------------------------------------------------
	*				    DATABASE CONNECTION
	*-----------------------------------------------------
	*Create a sql.DB abstraction with the given credentials
	 */

	db, err := sql.Open("mysql",
		DB_USERNAME+":"+DB_PASSWORD+"@tcp("+DB_HOST+":"+DB_PORT+")/"+DB_DATABASE)

	if err != nil {
		log.Fatal(err)
	}
	return db
}

var DB_CONN = OpenConn()
