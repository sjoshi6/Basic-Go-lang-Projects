package manager

import (
	"go-keystore/database/postgres"
	"hash/fnv"
	"log"
)

// StartManager : Start a manager to manage replication
func StartManager() {

	err := db.CreateIndexTableIfNotExists()
	log.Fatal(err)

	/*
	   Start a mux router here to accept new nodes and add their ids to db

	*/

}

func getIndexFromHash(listLength uint32, key string) uint32 {

	hashValue := getHashValue(key)
	index := hashValue % listLength
	return index
}

// Get has value from a key.
func getHashValue(key string) uint32 {

	hash := fnv.New32a()
	hash.Write([]byte(key))
	return hash.Sum32()

}
