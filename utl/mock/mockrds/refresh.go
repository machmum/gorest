package mockrds

import (
	"github.com/go-redis/redis"
)

// Refresh database mock
type Refresh struct {
	SaveFn   func(rdb *redis.Client, data interface{}, key, index string, lifetime uint) error
	DeleteFn func(rdb *redis.Client, key, index string) error
}

// FindByUsername mock
func (r *Refresh) Delete(rdb *redis.Client, key, index string) error {
	return r.DeleteFn(rdb, key, index)
}

// FindByUsername mock
func (r *Refresh) Save(rdb *redis.Client, data interface{}, key, index string, lifetime uint) error {
	return r.SaveFn(rdb, data, key, index, lifetime)
}
