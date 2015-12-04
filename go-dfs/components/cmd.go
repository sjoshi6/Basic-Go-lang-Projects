package components

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go-dfs/util"
	"net"
	"os"
	"strings"

	"github.com/garyburd/redigo/redis"
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

// CdHandler :  handler for cd function
func CdHandler(conn redis.Conn, cmdtype string, currDir string, dirName string) string {

	switch cmdtype {
	case "basedir":
		return fixedBaseDir

	case "back":

		if currDir == "/" {

			fmt.Println("Already at root.")
			return currDir

		}

		lastSlashIndx := strings.LastIndex(currDir, "/")
		secondlastIndx := strings.LastIndex(currDir[:lastSlashIndx], "/")
		currDir = currDir[:secondlastIndx+1]
		fmt.Printf("Moved to dir: %s \n", currDir)
		return currDir

	case "normal":

		intendedDir := currDir + dirName + "/"
		files, _ := redis.Strings(conn.Do("KEYS", "*"))
		found := false
		for _, file := range files {

			if strings.Index(file, intendedDir) == 0 {

				currDir = intendedDir
				fmt.Printf("Moved to dir: %s \n", currDir)
				found = true
				dirName = ""
				break
			}
		}

		if found == false {

			fmt.Println("Requested directory could not be found")
		}
		return currDir

	default:
		return currDir
	}

}

// LSHandler : Handler for the LS command
func LSHandler(conn redis.Conn, currDir string) {

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

				if !util.SliceContains(dirArr, path) {
					dirArr = append(dirArr, path)
				}

			}

		}
	}

	// Print the dirs found
	for _, dirname := range dirArr {

		fmt.Printf("dir \t--\t %s\n", "/"+dirname)
	}

	//Print the files found
	for _, filename := range fileArr {

		fmt.Printf("file \t--\t %s\n", "/"+filename)
	}

}

// HandleCat : Used to handle file cat commands
func HandleCat(cmdFile string) string {

	// Connect to master to get file Info
	mastertcpconn, err := net.Dial("tcp", ":5000")
	if err != nil {
		fmt.Print("No master tcpcon found")
	}

	// Send to master file server
	fmt.Fprintf(mastertcpconn, cmdFile+"\n")
	masterReply, _ := bufio.NewReader(mastertcpconn).ReadString('\n')

	// Unbox the received marshalled JSON
	var revIndex util.ReverseIndex
	unmarshallErr := json.Unmarshal([]byte(masterReply), &revIndex)

	if unmarshallErr != nil {
		fmt.Println(unmarshallErr)
	}

	// Print the received object
	fmt.Printf("\n%+v", revIndex)

	// Extract destination and full path
	destination := revIndex.Destination
	fullPath := revIndex.AbsolutePath

	// After Slave IP is found use this
	fmt.Println("Slave IP is found " + destination)

	// Connect to this slave to get cat data
	tcpconn, er := net.Dial("tcp", destination+":8080")
	if er != nil {
		fmt.Print("No tcpcon found")
	}

	// Read content from slave
	cmd := "cat " + fullPath + "\n"
	fmt.Fprintf(tcpconn, cmd)

	// Receive text from server
	message, _ := bufio.NewReader(tcpconn).ReadBytes('\t')
	return string(message)

}
