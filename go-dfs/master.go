package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
)

// ReceiveMessages : Receive messages fron master_messages redis channel
func ReceiveMessages(newSlaveChannel chan string, ipAddress string) {

	conn, err := redisurl.ConnectToURL(redisURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Close only when function exits
	defer conn.Close()

	// Creating a pubsubConn for master messages
	pubsubConn := redis.PubSubConn{Conn: conn}
	pubsubConn.Subscribe(masterMessageQueue)

	for {
		switch val := pubsubConn.Receive().(type) {

		case redis.Message:
			// If the data being received is a text message then push it to the channel
			newSlaveChannel <- string(val.Data)

		case redis.Subscription:
			//Handle Subscription here

		case error:
			return
		}

	}
}

// HandleNewSlaves : handles new slave start dir structure
func HandleNewSlaves(newSlaveChannel chan string) {

	for {
		newSlave := <-newSlaveChannel

		jsonObj := []byte(newSlave)
		var masterMsg MasterMessage

		err := json.Unmarshal(jsonObj, &masterMsg)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%+v", masterMsg)
		CreateFileMapping(masterMsg)

	}
}

//CreateFileMapping : Creates a mapping of file addresses
func CreateFileMapping(masterMsg MasterMessage) {
	conn, err := redisurl.ConnectToURL(redisURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var revIndex ReverseIndex

	for _, filepath := range masterMsg.FilePaths {
		index := strings.Index(filepath, "shared")
		filestart := index + 6
		relPath := filepath[filestart:]
		revIndex.AbsolutePath = filepath
		revIndex.Destination = masterMsg.IpAddress

		jsonObj, err := json.Marshal(revIndex)
		if err != nil {
			fmt.Println("Unable to marshal json")
		}

		// Insert the reverse Index in Redis
		conn.Do("SET", relPath, string(jsonObj))
	}

	conn.Close()
}

//GetFileIPServer :  Return File Destination IP and Absolute Path
func GetFileIPServer() {

	fmt.Println("Starting master IP getter server")
	ln, err := net.Listen("tcp", ":5000")

	if err != nil {

		fmt.Println(err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {

			fmt.Println(err)
		}

		go handleMasterConnection(conn)
	}

}

func handleMasterConnection(conn net.Conn) {

	message, _ := bufio.NewReader(conn).ReadString('\n')

	msg := message[:len(message)-1]
	// output message received
	fmt.Print("Command Received:", string(msg))

	// Extract filename and read its content
	relPath := string(msg)

	redisconn, err := redisurl.ConnectToURL(redisURL)
	if err != nil {

		fmt.Println(err)
	}

	defer redisconn.Close()

	// JSON Obj Values are received as bytes and translated to structs
	val, _ := redis.Bytes(redisconn.Do("GET", relPath))

	var revIndex ReverseIndex
	unmarshallErr := json.Unmarshal(val, &revIndex)

	if unmarshallErr != nil {
		fmt.Println(unmarshallErr)
	}

	fmt.Printf("\n %+v \n", revIndex)

	bytesJSON, _ := json.Marshal(revIndex)
	strJSON := string(bytesJSON) + "\n"
	// Send contents to client
	conn.Write([]byte(strJSON))
	conn.Close()

}
