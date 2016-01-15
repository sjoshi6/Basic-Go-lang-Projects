package rpcserv

import "sync"

// NewRPC : Factory for RPC objects
func NewRPC() (r *RPC) {
	return &RPC{
		Mu:       &sync.RWMutex{},
		Requests: &Requests{},
	}
}
