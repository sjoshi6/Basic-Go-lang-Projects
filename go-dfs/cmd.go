package main

import (
	"bufio"
	"fmt"
	"os"
)

// CommandLineInput : Take Inputs from command line
func CommandLineInput(commandChan chan string, exit *bool) {

	// Create a command line reader
	fmt.Println("Command Line Input active...")

	bufferedIO := bufio.NewReader(os.Stdin)
	for !*exit {

		line, _, err := bufferedIO.ReadLine()
		if err != nil {
			fmt.Print("error in input function")
		}

		cmdLine := string(line)

		if cmdLine == "/exit" {
			commandChan <- cmdLine
			break

		} else if cmdLine == "/dir" {

			commandChan <- cmdLine

		} else {

			fmt.Println("Unrecognized command")

		}

	}

	return

}

// CmdHandler : Handles all commands
func CmdHandler(commandChan chan string, exit *bool) {

	for !*exit {

		cmd := <-commandChan
		fmt.Printf("Command Received: %s \n", cmd)

		switch cmd {

		// If command is to quit program
		case "/exit":

			fmt.Println("Setting exit to true")
			*exit = true
			return
		case "/dir":
			dirStruct := GetDirStructure()
			for _, path := range dirStruct {
				fmt.Println(path)
			}
		default:
			continue
		}

	}
}
