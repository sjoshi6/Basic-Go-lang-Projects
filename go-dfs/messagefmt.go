package main

// MasterMessage : Format for messages to the master
type MasterMessage struct {
	IpAddress string
	FilePaths []string
}
