package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var (
	client = &redisClient{}
)

type redisClient struct {
	c *redis.Client
}

type LoggerRedis struct {
	Code         string      `json:"code"`
	Timestamp    time.Time   `json:"timestamp"`
	Id           int         `json:"id"`
	Repositories string      `json:"repositories"`
	Column       int         `json:"column"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
}

// GetClient get the redis client
func InitializeRedis() *redisClient {
	c := redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379", //6379
		Password:     "",
		DB:           0,                 // Gunakan database default
		DialTimeout:  5 * time.Second,   // Timeout untuk membuat koneksi baru
		ReadTimeout:  30 * time.Second,  // Timeout untuk pembacaan data
		WriteTimeout: 30 * time.Second,  // Timeout untuk penulisan data
		PoolSize:     50,                // Ukuran pool koneksi (jumlah koneksi maksimum yang aktif)
		MinIdleConns: 10,                // Jumlah koneksi idle minimum
		IdleTimeout:  300 * time.Second, // Timeout koneksi idle
	})

	if err := c.Ping().Err(); err != nil {
		fmt.Println("Unable to connect to redis " + err.Error())
		return nil
	}
	client.c = c
	return client
}

func (client *redisClient) GetKey(key string, src interface{}) error {
	val, err := client.c.Get(key).Result()
	if err == redis.Nil || err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), &src)
	if err != nil {
		return err
	}
	return nil
}

// SetKey set key
func (client *redisClient) SetKey(key string, value interface{}, expiration time.Duration) error {
	cacheEntry, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = client.c.Set(key, cacheEntry, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client *redisClient) DeleteKey(key string) error {

	iter := client.c.Scan(0, key, 0).Iterator()
	for iter.Next() {
		err := client.c.Del(iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func (client *redisClient) SettexKey(key string, value interface{}, expiration time.Duration) error {
	cacheSettex, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = client.c.SetXX(key, cacheSettex, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (client *redisClient) ExpireKey(key string, expiration time.Duration) error {
	err := client.c.Expire(key, expiration).Err()
	if err != nil {
		return err
	}
	return nil

}

func (client *redisClient) FlushAll() error {
	err := client.c.FlushAll().Err()
	if err != nil {
		return err
	}
	return nil
}
