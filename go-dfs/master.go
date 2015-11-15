package main

import (
	"fmt"
	"os"

	"github.com/soveran/redisurl"
)

// Master : Contains go code for master
func Master() {

	conn, err := redisurl.ConnectToURL(redisURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Close only when function exits
	defer conn.Close()
}
