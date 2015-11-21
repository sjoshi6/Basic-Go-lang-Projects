package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
)

// Master : Contains go code for master
func Master(newSlaveChannel chan string, ipAddress string) {

	conn, err := redisurl.ConnectToURL(redisURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	go ReceiveMessages(newSlaveChannel, ipAddress)
	go HandleNewSlaves(newSlaveChannel)
	go GetFileIP(conn)
	// Close only when function exits
	defer conn.Close()
}

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
		filestart := index + 5
		relPath := filepath[filestart:]
		revIndex.AbsolutePath = filepath
		revIndex.Destination = masterMsg.IpAddress
		fmt.Println(relPath)
		jsonObj, err := json.Marshal(revIndex)
		if err != nil {
			fmt.Println("Unable to marshal json")
		}
		//fmt.Println(jsonObj)
		conn.Do("SET", relPath, string(jsonObj), "NX")
	}
	defer conn.Close()
}

//GetFileIP :  Return File Destination IP and Absolute Path
func GetFileIP(redisconn redis.Conn) {
	// will listen for message to process ending in newline (\n)
	fmt.Println("Starting master server")
	ln, err := net.Listen("tcp", ":5000")
	if err != nil {
		// handle error
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleMasterConnection(conn, redisconn)
	}

}

func handleMasterConnection(conn net.Conn, redisconn redis.Conn) {

	message, _ := bufio.NewReader(conn).ReadString('\n')

	// output message received
	fmt.Print("Command Received:", string(message))

	// Extract filename and read its content
	relPath := strings.Split(string(message), " ")[1]
	relPathcnt := utf8.RuneCountInString(relPath)
	//index := strings.LastIndex(absPath, "/")
	fileName := relPath[:relPathcnt-1]
	fmt.Println(fileName)
	jsonString, _ := redisconn.Do("GET", fileName)
	//jsonObj := []byte(jsonString)
	fmt.Println(jsonString)

	//var dat map[string]interface{}

	//json.Unmarshal(jsonString, &dat)
	//ipAdress := jsonString["Destination"].(string)
	// absPath := dat["AbsolutePath"].(string)
	// Send contents to client
	//conn.Write([]byte("ipAdress:" + ipAdress + "absPath:" + absPath + "\n"))
	conn.Close()

}
