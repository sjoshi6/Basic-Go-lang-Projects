package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
)

// Slave : Contains go code for master
func Slave(ipAddress string, slaveExit *bool) {

	conn, err := redisurl.ConnectToURL(redisURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Close only when function exits
	defer conn.Close()

	// Add the slave to redis list
	val, err := conn.Do("SADD", "online_slaves", ipAddress)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if val == nil {
		fmt.Println("Insert error")
		os.Exit(1)
	}

	dirStruct := GetDirStructure()
	masterMsg := MasterMessage{IpAddress: ipAddress, FilePaths: dirStruct}

	jsonObj, err := json.Marshal(masterMsg)
	conn.Do("PUBLISH", masterMessageQueue, jsonObj)

}

// RegisterSlave : Used to register slave to redis
func RegisterSlave(conn redis.Conn, key string, value string) {

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

//GetDirStructure : Used to get the directory structure of the shared folder
func GetDirStructure() []string {

	currDir, err := os.Getwd()
	searchDir := string(currDir) + "/../shared/"

	fileList := []string{}
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)

		}
		return nil
	})

	if err != nil {
		fmt.Println("Could not execute find command")
		fmt.Println(err)
		os.Exit(1)
	}

	return fileList
}
