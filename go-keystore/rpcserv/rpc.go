package rpcserv

import (
	"go-keystore/database/postgres"
	"go-keystore/model"
	"log"
	"sync"
)

// RPC : A structure for RPC objects
type RPC struct {
	Requests *Requests
	Mu       *sync.RWMutex
}

// Requests : Used by RPC to track how many calls were made to each function
type Requests struct {
	Get    uint64
	Put    uint64
	Delete uint64
}

// Get : Used to get a keypair from storage node
func (t *RPC) Get(key string, keypair *model.KeyPair) error {

	// Get a read lock on the object
	t.Mu.RLock()
	defer t.Mu.RUnlock()

	value, _ := db.GetFromLocalNode(key)

	// Set the incoming second params value to a new KeyPair object with value from the DB
	*keypair = model.KeyPair{
		Key:   key,
		Value: string(value),
	}

	// Increment the number of get requests to the rpc obj
	t.Requests.Get++

	return nil
}

// Put : Used to put a keypair into a storagenodes local store
func (t *RPC) Put(keypair *model.KeyPair, success *bool) error {

	// Recieve a RW lock
	t.Mu.Lock()
	defer t.Mu.Unlock()

	// Increment request counter
	t.Requests.Put++

	key := keypair.Key
	value := keypair.Value

	err := db.InsertIntoLocalNode(key, value)
	if err != nil {

		log.Println("Error in inserting data to local node")

		// Set the success object to false
		*success = false
		return err
	}

	// If no err then set success to true
	*success = true
	return nil
}

// Delete : Used to delete a keypair from storagenode
func (t *RPC) Delete(key string, success *bool) error {

	// Receive a WR lock
	t.Mu.Lock()
	defer t.Mu.Unlock()

	// Increment requests counter
	t.Requests.Delete++

	err := db.DeleteFromLocalNode(key)
	if err != nil {
		*success = false
		return err
	}

	// Set the addressed object to true
	*success = true

	return nil
}
