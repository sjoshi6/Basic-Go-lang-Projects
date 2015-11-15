package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
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
	}
}
