package main

import (
	"fmt"
	"os"

	"github.com/soveran/redisurl"
)

func Slave() {

	managerConn, err := redisurl.ConnectToURL(redisURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
