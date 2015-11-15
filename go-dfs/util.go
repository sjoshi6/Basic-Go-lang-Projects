package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

// GetIPAddress : Used to get IP of the running machine
func GetIPAddress() string {
	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String()
			// process IP address
		}
	}

	return ""
}

//SendHeartBeat : Keeps user session alive
func SendHeartBeat(conn redis.Conn, key string, val string, slaveExit *bool) {

	// Continue sending heartbeats until chat exit is not called
	for !*slaveExit {
		// Here user is expected to exist
		_, err := conn.Do("SET", key, val, "XX", "EX", "100")

		if err != nil {
			fmt.Println(err)
			fmt.Println("Heart beat not accecpted")
			os.Exit(1)
		}

		// Send heartbeat every 80 seconds to avoid expiry
		time.Sleep(80 * time.Second)

	}

}
