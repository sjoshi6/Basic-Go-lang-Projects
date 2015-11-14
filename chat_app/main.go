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

	//See carefully the function is defined and run
	//Kind of like anonymous func in Java
	go func() {

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
				subscribeToRedisChan <- string(val.Data)

			case redis.Subscription:
				//Handle Subscription here

			case error:
				return
			}
		}

	}()

	// Send hearbeats to keep user alive
	go func() {

		// Continue sending heartbeats until chat exit is not called
		for !chatExit {

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

	}()

	// Put data to messages queue on redis DB
	go func() {

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
				chatExit = true
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
	}()

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
