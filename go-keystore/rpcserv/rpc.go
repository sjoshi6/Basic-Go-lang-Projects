package rpcserv

import (
	"go-keystore/database/postgres"
	"go-keystore/model"
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

	value, _ := db.GetJSONFromLocalNode(key)

	// Set the incoming second params value to a new KeyPair object with value from the DB
	*keypair = model.KeyPair{
		Key:   key,
		Value: string(value),
	}

	// Increment the number of get requests to the rpc obj
	t.Requests.Get++

	return nil
}
