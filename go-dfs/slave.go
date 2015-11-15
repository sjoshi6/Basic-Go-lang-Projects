package main

import (
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
)

// Slave : Contains go code for master
func Slave() {

	conn, err := redisurl.ConnectToURL(redisURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Close only when function exits
	defer conn.Close()

}

// RegisterSlave : Used to register slave to redis
func RegisterSlave(conn redis.Conn, ipAddress string) {

	key := "online." + ipAddress
	value := ipAddress + ":8000"

	// Register slave to redis
	val, err := conn.Do("SET", key, value, "NX", "EX", "100")

	// If DB throws err on insert
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}

	// If the insert is not a success and fails without ok message
	if val == nil {

		fmt.Println("Could not insert, Key exists in DB")
		fmt.Println("Slave is already online")
		os.Exit(1)

	}
}
