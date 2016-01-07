package db

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
)

const (
	redisURL = "redis://localhost:6379"
)

// GetRedisConn : Get redis conn for go-lbapp
func GetRedisConn() (redis.Conn, error) {

	// Connecting to the local redis DB instance
	conn, err := redisurl.ConnectToURL(redisURL)
	if err != nil {
		log.Println("Could not connect to redis DB")
		log.Println(err)

	}

	return conn, err
}

// GetListValueByKey : Get a list of values from redis
func GetListValueByKey(key string) ([]string, error) {

	redisConn, connerr := GetRedisConn()

	if connerr != nil {
		return nil, connerr
	}

	value, err := redis.Strings(redisConn.Do("LRANGE", key, 0, -1))
	if err != nil {
		log.Println("Error in reading list from redis")
		return nil, err
	}

	return value, nil
}

// InsertInList : Insert a new User in Event List
func InsertInList(key string, val string) error {

	redisConn, connerr := GetRedisConn()

	if connerr != nil {
		return connerr
	}

	_, err := redisConn.Do("LINSERT", key, val)
	if err != nil {

		log.Println("Could not insert value in List")
		return err
	}

	return nil
}
