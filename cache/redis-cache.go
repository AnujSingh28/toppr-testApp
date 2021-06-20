package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"test-app/entity"
	"time"

	"github.com/go-redis/redis/v8"
)

type TestCache interface {
	SetQuestion(key string, value entity.Questions)
	GetQuestion(key string) entity.Questions
	SetTestID(key string, value uint64)
	GetTestID(key string) uint64
	SetScore(key string, value int)
	GetScore(key string) int
	UpdateScore(key string)
	FlushAll()
}

type redisCache struct {
	host    string
	db      int
	expires time.Duration
	ctx     context.Context
}

func NewRedisCache(host string, db int, exp time.Duration, ctx context.Context) TestCache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: exp,
		ctx:     ctx,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) SetQuestion(key string, value entity.Questions) {
	client := cache.getClient()

	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	err = client.Set(cache.ctx, key, json, cache.expires*time.Second).Err()
	if err != nil {
		panic(err)
	}
}

func (cache *redisCache) GetQuestion(key string) entity.Questions {
	client := cache.getClient()

	val, err := client.Get(cache.ctx, key).Result()
	if err != nil {
		panic(err)
	}
	question := entity.Questions{}
	err = json.Unmarshal([]byte(val), &question)
	if err != nil {
		panic(err)
	}
	return question
}

func (cache *redisCache) SetTestID(key string, value uint64) {
	client := cache.getClient()

	err := client.Set(cache.ctx, key, value, cache.expires*time.Second).Err()
	if err != nil {
		panic(err)
	}
}

func (cache *redisCache) GetTestID(key string) uint64 {
	client := cache.getClient()
	value, err := client.Get(cache.ctx, key).Result()
	if err != nil {
		panic(err)
	}
	val, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return uint64(val)
}

func (cache *redisCache) SetScore(key string, value int) {
	client := cache.getClient()

	err := client.Set(cache.ctx, key, value, cache.expires*time.Second).Err()
	if err != nil {
		panic(err)
	}
}

func (cache *redisCache) GetScore(key string) int {
	client := cache.getClient()
	value, err := client.Get(cache.ctx, key).Result()
	if err != nil {
		panic(err)
	}
	val, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return val
}

func (cache *redisCache) UpdateScore(key string) {
	client := cache.getClient()
	client.Incr(cache.ctx, key)
}

func (cache *redisCache) FlushAll() {
	client := cache.getClient()
	client.FlushAll(cache.ctx)
}
