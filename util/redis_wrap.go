package util

import (
	"time"

	"strconv"

	redis "gopkg.in/redis.v3"
)

type RedisClient interface {
	Get(key string) StringCmd
	Set(key string, value interface{}, expiration time.Duration) StatusCmd
	Erase(key string)
	Exists(key string) bool
	SetAdd(key string, val string) IntCmd
	SetCard(key string) IntCmd
	SetList(key string) ([]string, error)
	SetRemove(key string, val ...string) IntCmd
	SortedSetAdd(key string, score int, value string) bool
	SortedSetRange(key string, min int64, max int64) ([]string, error)
	SortedSetRemoveRange(key string, min int64, max int64) IntCmd
	SortedSetRank(key string, value string) IntCmd
	SortedSetRemove(key string, values ...string) IntCmd
	SortedSetCard(key string) IntCmd
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

func (r *RedisClientWrap) SetCard(key string) IntCmd {
	return r.R.SCard(key)
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

func (r *RedisClientWrap) SetRemove(key string, val ...string) IntCmd {
	return r.R.SRem(key, val...)
}

func (r *RedisClientWrap) SortedSetAdd(key string, score int, value string) bool {
	i, _ := r.R.ZAdd(key, redis.Z{Score: float64(score), Member: value}).Result()
	return i == 1
}

func (r *RedisClientWrap) SortedSetRange(key string, min int64, max int64) ([]string, error) {
	return r.R.ZRangeByScore(key, redis.ZRangeByScore{Min: strconv.FormatInt(min, 10), Max: strconv.FormatInt(max, 10)}).Result()
}

func (r *RedisClientWrap) SortedSetRemoveRange(key string, min int64, max int64) IntCmd {
	return r.R.ZRemRangeByScore(key, strconv.FormatInt(min, 10), strconv.FormatInt(max, 10))
}

func (r *RedisClientWrap) SortedSetRank(key string, value string) IntCmd {
	return r.R.ZRank(key, value)
}

func (r *RedisClientWrap) SortedSetRemove(key string, values ...string) IntCmd {
	return r.R.ZRem(key, values...)
}

func (r *RedisClientWrap) SortedSetCard(key string) IntCmd {
	return r.R.ZCard(key)
}
