package db

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
)

// RedisConn : Get redis conn for go-lbapp
func RedisConn(URL string) (redis.Conn, error) {

	// Connecting to the local redis DB instance
	conn, err := redisurl.ConnectToURL(URL)
	if err != nil {
		log.Println("Could not connect to redis DB")
		log.Println(err)

	}

	return conn, err
}
