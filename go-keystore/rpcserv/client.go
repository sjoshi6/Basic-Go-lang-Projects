package rpcserv

import (
	"go-keystore/config"
	"go-keystore/model"
	"net"
	"net/rpc"
	"time"
)

// Client : A client object for rpc calls
type Client struct {
	connection *rpc.Client
}

// NewClient : Factory for RPC Clients
func NewClient(hostname string) (*Client, error) {

	// A complete server address to connect via tcp
	rpcServer := hostname + ":" + settings.ServerPort

	conn, err := net.DialTimeout("tcp", rpcServer, time.Millisecond*500)

	if err != nil {
		return nil, err
	}

	return &Client{connection: rpc.NewClient(conn)}, nil
}

// Get : A get call client wrapper
func (c *Client) Get(key string) (*model.KeyPair, error) {

	// Create a pointer to a keypair
	var keypair *model.KeyPair

	// Call the Get function and
	err := c.connection.Call("RPC.Get", key, &keypair)

	if err != nil {
		return nil, err
	}

	return keypair, nil
}
