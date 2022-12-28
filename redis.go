package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
)

const cacheExpirationFallback = 10 * time.Second

var cacheExpiration time.Duration

type RedisClient struct {
	ctx    context.Context
	client *redis.Client
}

func NewRedisClient(ctx context.Context) (*RedisClient, error) {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	expiration, err := strconv.Atoi(os.Getenv("CACHE_EXPIRATION_IN_SECS"))
	if err != nil {
		cacheExpiration = cacheExpirationFallback
	} else {
		cacheExpiration = time.Duration(expiration) * time.Second
	}

	db, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		db = 0
	}
	username := os.Getenv("REDIS_USERNAME")
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		Username:    username,
		Password:    password,
		DB:          db,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 3 * time.Second,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{
		ctx:    ctx,
		client: client,
	}, nil
}

func (r *RedisClient) Close() {
	err := r.client.Close()
	check(err)
}

func (r *RedisClient) Get(key string) (*TitleBasics, error) {
	t := &TitleBasics{}
	value, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, CacheMissError
		}
		return nil, err
	}

	err = json.Unmarshal(value, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *RedisClient) Set(key string, value *TitleBasics) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	cmd := r.client.Set(r.ctx, key, b, cacheExpiration)
	return cmd.Err()
}
