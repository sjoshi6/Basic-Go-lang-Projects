package db

import (
	"database/sql"
	"fmt"

	// Ued for connecting to postgres server
	_ "github.com/lib/pq"
)

// GetDBConn : Used for DB connections
func GetDBConn(dbname string) *sql.DB {

	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		fmt.Println("Error in connecting to PostGres DB")
		panic(err.Error())
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return db

}
