package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/soveran/redisurl"
)

const (
	redisURL = "redis://localhost:6379"
)

func main() {

	// Set default value of chat exit to false
	chatExit := false

	args := os.Args
	if len(args) != 2 {

		fmt.Println("Usage: chat_client.go <username>")
		os.Exit(1)
	}

	username := args[1]
	userDBKey := "online-" + username

	// Connecting to the local redis DB instance
	userManageConn, err := redisurl.ConnectToURL(redisURL)
	if err != nil {
		fmt.Println("Could not connect to redis DB")
		fmt.Println(err)
		os.Exit(1)
	}

	// Close the connection when program exits
	defer userManageConn.Close()

	// Insert Data into redis DB, if insert is a success then publish message
	val, err := userManageConn.Do("SET", userDBKey, username, "NX", "EX", "100")

	// If DB throws err on insert
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}

	// If the insert is not a success and fails without ok message
	if val == nil {

		fmt.Println("Could not insert, Key exists in DB")
		fmt.Println("User is online on another terminal...")
		os.Exit(1)

	} else {

		// If user key was added , try to insert insert user to set of online users

		setAddVal, err := userManageConn.Do("SADD", "users", username)

		if err != nil {
			fmt.Println(err)
		}
		if setAddVal == nil {
			fmt.Println("User is online on another terminal...")
		}

		// If user was added to the set of online users then publish a new message

		userManageConn.Do("PUBLISH", "messages", "New User Joined: "+username)
		fmt.Println("You joined the application successfully.")
		fmt.Println("Chat Ready. Just Start typing.....")
	}

	// New channel for each user connection to read messages from redisDB
	subscribeToRedisChan := make(chan string)

	// Call to handle Subscribe of redis
	go SubscribeHandler(subscribeToRedisChan)

	// Send hearbeats to keep user alive
	go SendHeartBeat(userManageConn, username, userDBKey, &chatExit)

	// Handle publishing and commands
	go PublishAndCommandHandler(username, &chatExit)

	// While chatExit is not called poll messages from the subscribeChannel
	// This is why we defined the subscribeChannel outside function
	for !chatExit {
		select {
		case line := <-subscribeToRedisChan:

			// Read only messages from others
			if strings.Index(line, username) != 0 {
				fmt.Printf("%s \n", line)
			}

		default:
			// Sleep a second before message arrives again
			time.Sleep(100 * time.Millisecond)

		}
	}

	// Before the function exits remove the user from redis DB
	userManageConn.Do("DEL", userDBKey)
	userManageConn.Do("SREM", "users", username)
	userManageConn.Do("PUBLISH", "messages", "Sad :( User Left :- \t"+username)

}
