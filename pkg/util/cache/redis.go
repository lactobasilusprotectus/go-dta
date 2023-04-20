package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-dts/pkg/common/config"
	"log"
	"time"
)

var (
	ErrNilReturned = redis.ErrNil
)

type RedisClient struct {
	pool redis.Pool
}

func NewRedisClient(config config.RedisConfig) *RedisClient {
	return &RedisClient{
		pool: redis.Pool{
			MaxIdle:     3,
			IdleTimeout: time.Duration(500) * time.Millisecond,
			MaxActive:   100,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", config.Host)
				if err != nil {
					err = fmt.Errorf("redis: %+v\n", err)
					return nil, err
				}

				return conn, nil
			},
		}}
}

// Get gets setValue of the given setKey from Redis and returns a reply.
// You need to convert the reply using the appropriate function.
// https://godoc.org/github.com/gomodule/redigo/redis#hdr-Executing_Commands
func (r *RedisClient) Get(key string) (reply interface{}, err error) {
	conn := r.pool.Get() // get from pool
	defer func() {
		errClose := conn.Close() // don't forget to close
		if errClose != nil {
			log.Println("closing redis conn err:", errClose)
		}
	}()

	return conn.Do("GET", key)
}

// Set sets setKey-setValue pair that will expire after expireSeconds seconds.
// expireSeconds <= 0 means the setKey-setValue has no expiration.
func (r *RedisClient) Set(key string, value interface{}, expireSeconds int) (err error) {
	conn := r.pool.Get() // get from pool
	defer func() {
		errClose := conn.Close() // don't forget to close
		if errClose != nil {
			log.Println("closing redis conn err:", errClose)
		}
	}()

	if expireSeconds <= 0 {
		_, err = conn.Do("SET", key, value)
	} else {
		_, err = conn.Do("SETEX", key, expireSeconds, value)
	}
	return
}

// == CONVERTER FUNCTIONS ==
// other converter functions can be added based on the official docs of redigo:
// https://godoc.org/github.com/gomodule/redigo/redis#hdr-Executing_Commands
// See the unit tests for samples.

// String https://godoc.org/github.com/gomodule/redigo/redis#String
func String(reply interface{}, err error) (string, error) {
	return redis.String(reply, err)
}

// Int https://godoc.org/github.com/gomodule/redigo/redis#Int
func Int(reply interface{}, err error) (int, error) {
	return redis.Int(reply, err)
}

// Int64 https://godoc.org/github.com/gomodule/redigo/redis#Int64
func Int64(reply interface{}, err error) (int64, error) {
	return redis.Int64(reply, err)
}

// Bytes Use this for retrieving JSON object as bytes.
// https://godoc.org/github.com/gomodule/redigo/redis#Bytes
func Bytes(reply interface{}, err error) ([]byte, error) {
	return redis.Bytes(reply, err)
}
