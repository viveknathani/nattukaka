package cache

import (
	"github.com/gomodule/redigo/redis"
)

// Cache contains the connection pool for Redis.
// Server should call Initialize before usage and
// each goroutine should work on a new connection
// using Pool.Get().
type Cache struct {
	Pool *redis.Pool
}

// Initialize will set up the connection pool so that
// future connections can take place on the given URL.
func (cache *Cache) Initialize(url string, username string, password string) {

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {

			var conn redis.Conn
			var err error
			if username == "" || password == "" {
				conn, err = redis.Dial("tcp", url)
			} else {
				user := redis.DialUsername(username)
				pass := redis.DialPassword(password)
				conn, err = redis.Dial("tcp", url, user, pass)
			}

			if err == nil {
				return conn, nil
			}

			return nil, err
		},
	}
	cache.Pool = pool
}

// Close will free up all the resources of the pool.
func (cache *Cache) Close() error {
	return cache.Pool.Close()
}

// Get is a generic method to do the obvious.
// It is not a method of the Cache struct because connection management
// is server's responsibility.
func Get(conn redis.Conn, key string) ([]byte, error) {
	v, err := conn.Do("GET", key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, err
	}
	return v.([]byte), err
}

// Set is a generic method to do the obvious.
// It is not a method of the Cache struct because connection management
// is server's responsibility.
func Set(conn redis.Conn, key string, value []byte) (interface{}, error) {
	return conn.Do("SET", key, value)
}
