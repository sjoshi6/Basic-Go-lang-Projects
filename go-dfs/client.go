package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/soveran/redisurl"
)

const (
	fixedBaseDir = "/"
)

// FileSystemCommandHandler : Manages commands issued by the client
func FileSystemCommandHandler(exit *bool, username string) {

	// Very important var ; used by the commands for reference
	currDir := "/"

	conn, err := redisurl.ConnectToURL(redisURL)
	if err != nil {
		fmt.Println("could not connect to the redis URL")
	}

	for !*exit {
		// Create a command line reader
		fmt.Println("File System commands Available...")

		bufferedIO := bufio.NewReader(os.Stdin)
		for !*exit {

			fmt.Printf("%s-cmd-prompt $:  ", username)
			line, _, err := bufferedIO.ReadLine()
			if err != nil {
				fmt.Print("error in input function")
			}

			// command issued as a string
			cmdLine := string(line)

			dir := "" // Needs to be outside loop

			if cmdLine == "cd /" {

				//Pre handler for cd command
				cmdLine = "cd basedir"

			} else if cmdLine == "cd .." {

				// Prehandler for cd .. command
				cmdLine = "cd back"
			} else {

				// Handler for cd <dir> command

				rp := regexp.MustCompile("cd [a-z]+")
				iscdCmd := rp.MatchString(cmdLine)
				if iscdCmd == true {

					dir = strings.Split(cmdLine, " ")[1]
					cmdLine = "cd"
				}

			}

			switch cmdLine {

			case "ls":

				LSHandler(conn, currDir)
				break

			case "cd basedir":

				currDir = CdHandler(conn, "basedir", currDir, "")
				fmt.Printf("Moved to dir: %s \n", currDir)
				break

			case "cd back":

				currDir = CdHandler(conn, "back", currDir, "")

				break

			case "cd":

				currDir = CdHandler(conn, "normal", currDir, dir)
				dir = ""
				break

			case "exit":

				fmt.Println("Program exiting")
				*exit = true

			case "cat":
				fmt.Println("Executing cat command")
				tcpconn, er := net.Dial("tcp", "192.168.0.6:8080")
				if er != nil {
					fmt.Print("No tcpcon found")
				}
				// send to socket
				fmt.Fprintf(tcpconn, "cat test/sau.txt\n")

				// Receive text from server
				message, _ := bufio.NewReader(tcpconn).ReadString('\n')
				fmt.Println(message)

			default:
				fmt.Println("Unrecognized command")

			}

			fmt.Println()

		}

		return

	}
}
