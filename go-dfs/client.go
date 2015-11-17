package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/garyburd/redigo/redis"
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

				files, _ := redis.Strings(conn.Do("KEYS", "*"))
				baseDirEndIndex := len(currDir)
				dirArr := []string{}
				fileArr := []string{}

				for _, file := range files {
					if strings.Index(file, currDir) == 0 {
						relpath := file[baseDirEndIndex:]
						fileordir := strings.Index(relpath, "/")

						if fileordir == -1 {

							fileArr = append(fileArr, relpath)

						} else {

							path := relpath[:fileordir+1]

							if !SliceContains(dirArr, path) {
								dirArr = append(dirArr, path)
							}

						}

					}
				}

				// Print the dirs found
				for _, dirname := range dirArr {

					fmt.Printf("dir \t--\t %s\n", dirname)
				}

				//Print the files found
				for _, filename := range fileArr {

					fmt.Printf("file \t--\t %s\n", filename)
				}

				break

			case "cd basedir":

				currDir = fixedBaseDir
				fmt.Printf("Moved to dir: %s \n", currDir)
				break

			case "cd back":

				if currDir == "/" {

					fmt.Println("Already at root.")

				} else {

					lastSlashIndx := strings.LastIndex(currDir, "/")
					secondlastIndx := strings.LastIndex(currDir[:lastSlashIndx], "/")
					currDir = currDir[:secondlastIndx+1]
					fmt.Printf("Moved to dir: %s \n", currDir)
				}

				break

			case "cd":

				intendedDir := currDir + dir + "/"
				files, _ := redis.Strings(conn.Do("KEYS", "*"))
				found := false
				for _, file := range files {

					if strings.Index(file, intendedDir) == 0 {

						currDir = intendedDir
						fmt.Printf("Moved to dir: %s \n", currDir)
						found = true
						dir = ""
						break
					}
				}

				if found == false {

					fmt.Println("Requested directory could not be found")

				}
				break

			case "exit":

				fmt.Println("Program exiting")
				*exit = true

			default:
				fmt.Println("Unrecognized command")

			}

			fmt.Println()

		}

		return

	}
}
