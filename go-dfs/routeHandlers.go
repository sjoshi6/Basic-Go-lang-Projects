package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
)

// HomePage : Home Page route
func HomePage(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Welcome!")
}

// GetAll : All file paths printer
func GetAll(w http.ResponseWriter, r *http.Request) {

	currDir := "/"

	var conn, _ = redisurl.ConnectToURL(redisURL)

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

		fmt.Fprintf(w, "dir \t--\t %s\n", dirname)
	}

	//Print the files found
	for _, filename := range fileArr {

		fmt.Fprintf(w, "file \t--\t %s\n", filename)
	}

	conn.Close()

}
