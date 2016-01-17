package db

import (
	"database/sql"
	"fmt"
	"go-keystore/config"
	"log"

	// Used for connecting to postgres server
	_ "github.com/lib/pq"
)

// CreateDBIfNotExists : Create DB if not present
func CreateDBIfNotExists() error {

	// Connect to DB without dbname
	dbConnStr := fmt.Sprintf("root:root@tcp(%s:%s)?charset=utf8", settings.DBHostName, settings.DBPort)

	db, err := sql.Open("mysql", dbConnStr)
	defer db.Close()

	if err != nil {
		log.Println("Failed to create the DB")
		log.Println(err)
		return err
	}

	result, err := db.Exec("CREATE DATABASE IF NOT EXISTS $1", settings.DBName)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(result)

	return nil
}

// CreateTableIfNotExists : Create table if not exist
func CreateTableIfNotExists() error {

	// Create DB conn
	db, err := GetDBConn()
	if err != nil {
		log.Println("Error Connecting to DB")
		return err
	}

	// Defer db close
	defer db.Close()

	// Creating the table
	result, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS KeyPair ( key VARCHAR(200) PRIMARY KEY, value JSON NOT NULL);")

	if err != nil {
		return err
	}

	log.Println(result)

	return nil
}

// GetDBConn : conn object for DB - Make sure function closes it
func GetDBConn() (*sql.DB, error) {

	dbconnStr := fmt.Sprintf("dbname=%s sslmode=disable", settings.DBName)
	db, err := sql.Open("postgres", dbconnStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}
