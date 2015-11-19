package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
)

var fileDir, _ = os.Getwd()
var sharedDir = string(fileDir) + "/../shared/"

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

// StartFileServer : Used to start file server at slave
func StartFileServer() {

	fmt.Println("Started File Server")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	// will listen for message to process ending in newline (\n)
	message, _ := bufio.NewReader(conn).ReadString('\n')

	// output message received
	fmt.Print("Command Received:", string(message))

	// Extract filename and read its content
	absPath := strings.Split(string(message), " ")[1]
	absPathcnt := utf8.RuneCountInString(absPath)
	fileName := absPath[:absPathcnt-1]

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	// Send contents to client
	conn.Write([]byte(string(data) + "\n"))
	conn.Close()

}
