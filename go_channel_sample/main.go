package main

import (
	"fmt"
	"time"
)

func main() {

	messages := make(chan string)

	fmt.Println("Starting Server")
	// Print data to cmd line

	// get messages from cmd ling
	for i := 0; true; i++ {

		go printMsg(messages)
		go getMsgs(messages)

		//sleep for a while so that the program doesn’t exit immediately
		time.Sleep(2 * 1e9)

	}

	defer close(messages)
}

func getMsgs(messagesChan chan string) {

	fmt.Print("Enter text: ")
	var input string
	fmt.Scanln(&input)
	messagesChan <- input

}

func printMsg(messagesChan chan string) {

	fmt.Println("Waiting for a new message...")
	outputMsg := <-messagesChan
	fmt.Println("Ping from Channel", outputMsg)

}
