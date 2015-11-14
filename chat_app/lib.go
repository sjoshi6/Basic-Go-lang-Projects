package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
)

//SendHeartBeat : Keeps user session alive
func SendHeartBeat(userManageConn redis.Conn, username string, userDBKey string, chatExit *bool) {

	// Continue sending heartbeats until chat exit is not called
	for !*chatExit {

		// Here user is expected to exist
		_, err := userManageConn.Do("SET", userDBKey, username, "XX", "EX", "100")

		if err != nil {
			fmt.Println(err)
			fmt.Println("Heart beat not accecpted")
			os.Exit(1)
		}

		// Send heartbeat every 80 seconds to avoid expiry
		time.Sleep(80 * time.Second)

	}

}

// SubscribeHandler : Subsribes to redis to fetch messages
func SubscribeHandler(subChannel chan string) {
	// create a subscribe connection to RedisDB
	subscribeConn, err := redisurl.ConnectToURL("redis://localhost:6379")
	if err != nil {

		fmt.Println(err)
		os.Exit(1)

	}

	// Before function exits close the connection
	defer subscribeConn.Close()

	pubsubConn := redis.PubSubConn{Conn: subscribeConn}
	pubsubConn.Subscribe("messages") // Subscribed to messages list in redis DB

	for {

		switch val := pubsubConn.Receive().(type) {

		case redis.Message:
			// If the data being received is a text message then push it to the channel
			subChannel <- string(val.Data)

		case redis.Subscription:
			//Handle Subscription here

		case error:
			return
		}
	}

}

// PublishAndCommandHandler : function to publish and handle commands for the code
func PublishAndCommandHandler(username string, chatExit *bool) {

	publishConn, err := redisurl.ConnectToURL("redis://localhost:6379")
	if err != nil {

		fmt.Println(err)
		os.Exit(1)

	}

	// Before function exits close the connection
	defer publishConn.Close()

	// Create a command line reader
	bufferedIO := bufio.NewReader(os.Stdin)
	for {
		line, _, err := bufferedIO.ReadLine()
		if err != nil {
			fmt.Print("error in input function")
		}

		// If user enters an exit request quit it
		if string(line) == "/exit" {

			// If the command is to exit. Quit function by return
			*chatExit = true
			return

		} else if string(line) == "/online" {

			// Get all the names in string format from users
			names, _ := redis.Strings(publishConn.Do("SMEMBERS", "users"))

			// Range first param is index and second is value
			for _, name := range names {

				// Print list of Online users
				fmt.Printf("Online: %s \n", name)
			}

		} else if strings.Index(string(line), "/") == 0 {
			// If the string begins with / and no other command its mostly a typo
			// Ignore this command input
		} else {

			publishConn.Do("PUBLISH", "messages", username+":- \t"+string(line))
		}

	}
}
