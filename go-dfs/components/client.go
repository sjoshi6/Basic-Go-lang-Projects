package components

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/soveran/redisurl"
)

const (
	fixedBaseDir       = "/"
	redisURL           = "redis://152.46.16.250:6379"
	masterMessageQueue = "master_message"
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

			// --------------- Pre Hanlder for cat command ----------- //

			cmdFile := ""

			if strings.Index(cmdLine, "cat") == 0 {

				cmdFile = strings.Split(cmdLine, " ")[1]
				cmdLine = "cat"

			}

			//--------------- Pre handlers for cd commands ------------ //
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

				data := HandleCat(cmdFile)
				fmt.Println(data)

			default:
				fmt.Println("Unrecognized command")

			}

			fmt.Println()

		}

		return

	}
}
