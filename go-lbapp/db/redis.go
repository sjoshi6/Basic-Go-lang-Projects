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

// RedisGetListValueByKey : Get a list of values from redis
func RedisGetListValueByKey(key string) ([]string, error) {

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

// RedisInsertInSet : Insert a new User in Event List
func RedisInsertInSet(key string, val string) error {

	redisConn, connerr := GetRedisConn()

	if connerr != nil {
		return connerr
	}

	log.Printf("Attempting to add user : %s to eventid : %s ", val, key)

	reply, err := redisConn.Do("SADD", key, val)
	if err != nil {

		log.Println("Could not insert value in List")
		log.Println(err)

		return err
	}

	if reply == nil {

		// Dont throw error just log that second subscribe was sent by same user
		log.Println("User is already subscribed to the event.")
		return nil
	}

	return nil
}

// RedisRemoveFromSet : Insert a new User in Event List
func RedisRemoveFromSet(key string, val string) error {

	redisConn, connerr := GetRedisConn()

	if connerr != nil {
		return connerr
	}

	log.Printf("Attempting to remove user : %s to eventid : %s ", val, key)

	reply, err := redisConn.Do("SREM", key, val)
	if err != nil {

		log.Println("Could not remove value from the list")
		log.Println(err)

		return err
	}

	if reply == nil {

		// Dont throw error just log that second subscribe was sent by same user
		log.Println("Err: User could not be removed.")
		return nil
	}

	return nil
}

// RedisCheckDuplicateSubscribe : Confirms if the Subscriber isnt already present.
func RedisCheckDuplicateSubscribe(key string, val string) (bool, error) {

	redisConn, connerr := GetRedisConn()

	// If err in conn : reply not subscribed but has err
	if connerr != nil {
		return false, connerr
	}

	reply, err := redisConn.Do("SISMEMBER", key, val)

	if err != nil {

		// If err in conn : reply not subscribed but has err
		log.Println(err)
		return false, err
	}

	if reply.(int64) == 1 {
		// Already subscribed and no err
		return true, nil
	}

	// Not subscribed neither err
	return false, nil
}
