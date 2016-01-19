package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-keystore/config"
	"log"

	// Used for connecting to postgres server
	_ "github.com/lib/pq"
)

// CreateStorageTableIfNotExists : Create table if not exist
func CreateStorageTableIfNotExists() error {

	log.Println("Validating the presence of keypair table on storage node...")

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
	log.Println("Database and table keypair is ready as required.")

	return nil
}

// CreateIndexTableIfNotExists : Creates an index table if not present
func CreateIndexTableIfNotExists() error {

	log.Println("Validating the presence of index table on manager node...")

	// Create DB conn
	db, err := GetDBConn()
	if err != nil {
		log.Println("Error Connecting to DB")
		log.Println(err)
		return err
	}

	// Defer db close
	defer db.Close()

	// Creating the table
	result, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS ManagerIndex (id SERIAL PRIMARY KEY, IPAddress VARCHAR(20) NOT NULL);")

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(result)
	log.Println("Table index is ready as required.")

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

// GetFromLocalNode : JSON for value
func GetFromLocalNode(key string) ([]byte, error) {

	var value json.RawMessage

	log.Printf("Get value for %s \n", key)

	// Create DB conn
	db, err := GetDBConn()
	if err != nil {
		log.Println("Error Connecting to DB")
		return nil, err
	}

	// Defer db close
	defer db.Close()

	db.QueryRow("SELECT * FROM KeyPair where key=$1", key).Scan(&key, &value)

	return value, nil

}

// InsertIntoLocalNode : Insert a JSON into postgres
func InsertIntoLocalNode(key string, value string) error {

	// Create DB conn
	db, err := GetDBConn()
	if err != nil {
		log.Println("Error Connecting to DB")
		return err
	}

	// Defer db close
	defer db.Close()

	/*
			   Creating a very naive delete and insert implementation.
		       Later improve it with versioning of objects
	*/

	db.Exec("DELETE FROM KeyPair WHERE key=$1", key)
	_, inserterr := db.Exec("INSERT INTO KeyPair VALUES($1, $2);", key, value)

	if inserterr != nil {

		log.Println("Error in insert operation")
		log.Println(inserterr)

		return inserterr
	}

	return nil
}

// DeleteFromLocalNode : Delete a key from Local Node
func DeleteFromLocalNode(key string) error {

	// Create DB conn
	db, err := GetDBConn()
	if err != nil {

		log.Println("Error Connecting to DB")
		return err
	}

	// Defer db close
	defer db.Close()

	_, deleteerr := db.Exec("DELETE FROM KeyPair WHERE key=$1", key)
	if deleteerr != nil {

		return deleteerr
	}

	return nil
}
