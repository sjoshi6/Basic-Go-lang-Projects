package main

import (
	"fmt"
	"go-lbapp/api"
	"time"
)

func main() {

	fmt.Println("Go Server - Logs", time.Now())
	go api.StartServer()

}
