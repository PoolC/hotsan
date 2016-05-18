package util

import (
	"time"

	redis "gopkg.in/redis.v3"
)

type RedisClient interface {
	Get(key string) StringCmd
	Set(key string, value interface{}, expiration time.Duration) StatusCmd
	Erase(key string)
	Exists(key string) bool
	SetAdd(key string, val string) IntCmd
	SetList(key string) ([]string, error)
	Keys(pattern string) ([]string, error)
}

type Command interface {
	Val() string
	Result() (string, error)
	String() string
}

type StringCmd interface {
	Command
	Bytes() ([]byte, error)
	Float64() (float64, error)
	Int64() (int64, error)
	Scan(val interface{}) error
	Uint64() (uint64, error)
}

type IntCmd interface {
	Val() int64
	String() string
	Result() (int64, error)
}

type StatusCmd interface {
	Command
}

type RedisClientWrap struct {
	R *redis.Client
}

func (r *RedisClientWrap) Get(key string) StringCmd {
	return r.R.Get(key)
}

func (r *RedisClientWrap) Set(key string, value interface{}, expiration time.Duration) StatusCmd {
	return r.R.Set(key, value, expiration)
}

func (r *RedisClientWrap) Exists(key string) bool {
	return r.R.Exists(key).Val()
}

func (r *RedisClientWrap) SetAdd(key string, val string) IntCmd {
	return r.R.SAdd(key, val)
}

func (r *RedisClientWrap) SetList(key string) ([]string, error) {
	return r.R.SMembers(key).Result()
}

func (r *RedisClientWrap) Keys(pattern string) ([]string, error) {
	return r.R.Keys(pattern).Result()
}

func (r *RedisClientWrap) Erase(key string) {
	r.R.Del(key)
}
