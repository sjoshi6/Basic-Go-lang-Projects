package main

// ReverseIndex : Used to store Reverse Index information in Redis
// Useful for marshal and unmarshalling
type ReverseIndex struct {
	PrefixPath  string
	Destination string
}
