package cache

import (
	"context"
	"errors"
	"sync"

	"github.com/MojtabaArezoomand/lru_cache/internal/config"
	linkedlist "github.com/MojtabaArezoomand/lru_cache/internal/linked_list"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Cache is the LRU cache struct
	Cache struct {
		m        sync.Mutex
		list     *linkedlist.DoublyLinkedList
		storage  map[string]*linkedlist.Node
		capacity uint64
	}

	// getResult is the struct for sending cache's Get result using channels.
	getResult struct {
		val any
		err error
	}
)

// Errors.
var (
	ErrNotFound error = errors.New("not found")
)

// NewCache returns a new cache.
func NewCache() *Cache {
	var cfg config.CacheConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	cache := Cache{
		list:     linkedlist.NewDoublyLinkedList(),
		storage:  make(map[string]*linkedlist.Node),
		capacity: cfg.CacheCapacity.ToUint64(),
	}

	return &cache
}

// Get fetches the key from the cache.
func (c *Cache) Get(ctx context.Context, key string) (any, error) {
	if ctx == nil {
		panic("Context cannot be nil.")
	}

	getChan := make(chan getResult, 1)

	go func() {
		val, err := c.get(key)
		getChan <- getResult{val: val, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-getChan:
		return res.val, res.err
	}
}

// get fetches the key from storage.
func (c *Cache) get(key string) (any, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if node, ok := c.storage[key]; !ok {
		return nil, ErrNotFound
	} else {
		c.list.MoveToBack(node)
		val := node.GetVal()
		return val, nil
	}
}

// Set sets or overwrites the key-value to cache.
func (c *Cache) Set(ctx context.Context, key string, val any) error {
	if ctx == nil {
		panic("Context cannot be nil.")
	}

	done := make(chan bool, 1)

	go func() {
		c.set(key, val)
		done <- true
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

// set sets or overwrites the key-value to cache.
func (c *Cache) set(key string, val any) {
	c.m.Lock()
	defer c.m.Unlock()

	if node, ok := c.storage[key]; ok {
		node.SetVal(val)
		c.list.MoveToBack(node)
	} else {
		if c.capacity == c.list.Size() {
			key := c.list.RemoveHead()
			delete(c.storage, key)
		}

		node := c.list.AddToBack(key, val)
		c.storage[key] = node
	}
}
