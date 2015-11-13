package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Used fo connecting to postgres
const (
	DbName = "go_postgres_db"
)

func main() {

	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", DbName)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		fmt.Println("Error in connecting to MySQL DB")
		panic(err.Error())
	}

	// Closed DB connection after the connect function exits
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Select statement
	rows, err := db.Query("SELECT * FROM users;")
	for rows.Next() {
		var userid string
		var password string

		err = rows.Scan(&userid, &password)
		fmt.Println("User Id : Password")
		fmt.Printf(" %s : %s \n", userid, password)
	}

}
