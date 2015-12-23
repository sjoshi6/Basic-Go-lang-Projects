package main

import (
	"fmt"
	"go-lbapp/api"
	"go-lbapp/generics"
	"time"
)

// Used for connecting to postgres

func main() {

	controller := make(chan generics.SyncMsg)
	fmt.Println("Go Server - Logs", time.Now())

	go api.StartServer(controller)
	<-controller
}
