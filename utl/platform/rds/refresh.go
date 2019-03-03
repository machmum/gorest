package rds

import (
	"time"
	"github.com/go-redis/redis"
	"encoding/json"
)

// NewProfile returns a new profile redis instance
func NewRefresh() *Refresh {
	return &Refresh{}
}

// User represents the client for user table
type Refresh struct{}

// Find key in redis
func (r *Refresh) Find(rdb *redis.Client, key, index string) (string, error) {
	var err error

	data, err := rdb.HGetAll(key).Result()

	return data[index], err
}

// Save to key redis
// key data: index - data
func (r *Refresh) Save(rdb *redis.Client, data interface{}, key, index string, lifetime uint) error {
	var err error

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// write redis
	if _, err = rdb.HMSet(key, map[string]interface{}{index: bytes}).Result(); err != nil {
		return err
	}

	// set expire redis
	if _, err = rdb.Expire(key, time.Duration(lifetime)*time.Second).Result(); err != nil {
		return err
	}

	return err
}

// Delete key in redis
func (r *Refresh) Delete(rdb *redis.Client, key, index string) error {
	var err error

	_, err = rdb.HDel(key, index).Result()

	return err
}
