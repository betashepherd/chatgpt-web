package cache

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var RedisPool *redis.Pool

func Init() {
	network := ""
	address := ""
	password := ""
	rdbindex := 0
	maxIdle := 60
	RedisPool = &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: 120 * time.Second,
		Dial: func() (redis.Conn, error) {
			cli, err := redis.Dial(network, address)
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := cli.Do("AUTH", password); err != nil {
					cli.Close()
					return nil, err
				}
			}

			if rdbindex > 0 {
				if _, err := cli.Do("SELECT", rdbindex); err != nil {
					cli.Close()
					return nil, err
				}
			}

			return cli, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

}
