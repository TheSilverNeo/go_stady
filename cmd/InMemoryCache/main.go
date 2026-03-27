package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"
	"time"
)

type Cache interface {
	Set(key string, value any, ttl time.Duration)
	Get(key string) any
}

type DataWithTTL struct {
	value  any
	ttlCtx context.Context
}

type InMemoryCache struct {
	Shard []Shard
}

func (i *InMemoryCache) Set(key string, value any, ttl time.Duration) error {
	hash, err := hashKey(key)
	if err != nil {
		return err
	}

	shardId := *hash % uint64(len(i.Shard))
	i.Shard[shardId].Set(key, value, ttl)

	return nil
}

func (i *InMemoryCache) Get(key string) any {
	hash, err := hashKey(key)
	if err != nil {
		return err
	}

	shardId := *hash % uint64(len(i.Shard))
	return i.Shard[shardId].Get(key)
}

type Shard struct {
	data map[any]DataWithTTL
	mu   sync.RWMutex
}

func (s *Shard) Set(key string, value any, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), ttl*time.Second)
	defer cancel()
	data := DataWithTTL{value: value, ttlCtx: ctx}
	s.data[key] = data
}

func (s *Shard) Get(key string) any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if val, ok := s.data[key]; ok {
		return val.value
	}

	return nil
}

func newInMemoryCache(numShards int) *InMemoryCache {
	shards := make([]Shard, 0, numShards)
	for i := 0; i < numShards; i++ {
		shards = append(shards, Shard{data: make(map[any]DataWithTTL)})
	}

	return &InMemoryCache{shards}
}

func hashKey(key string) (*uint64, error) {
	hasher := fnv.New64a()

	_, err := hasher.Write([]byte(key))
	if err != nil {
		return nil, err
	}

	res := hasher.Sum64()

	return &res, nil
}

func main() {
	cache := newInMemoryCache(5)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := cache.Set("foo", "first", 5*time.Second)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := cache.Set("bar", "second", 5*time.Second)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	wg.Wait()

	fmt.Println(cache.Get("foo"))
}
