package cache

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCache(t *testing.T) {
	cache := NewCache()

	assert.NotNil(t, cache.list)
	assert.NotNil(t, cache.storage)
	assert.EqualValues(t, 2048, cache.capacity)

	os.Setenv("CACHE_CAPACITY", "0")

	assert.Panics(t, func() {
		NewCache()
	})

	os.Setenv("CACHE_CAPACITY", "2048")
}

func TestGetSet(t *testing.T) {
	cache := NewCache()
	cache.capacity = 3

	cache.set("first", 1)

	assert.EqualValues(t, 1, cache.list.Size())
	assert.EqualValues(t, 1, len(cache.storage))

	res, err := cache.get("first")
	assert.NoError(t, err)

	assert.Equal(t, 1, res)

	// Overwriting the first key
	cache.set("first", 2)

	assert.EqualValues(t, 1, cache.list.Size())
	assert.EqualValues(t, 1, len(cache.storage))

	res, err = cache.get("first")
	assert.NoError(t, err)

	assert.Equal(t, 2, res)

	cache.set("second", 4)
	cache.set("third", 5)

	assert.EqualValues(t, 3, cache.list.Size())
	assert.EqualValues(t, 3, len(cache.storage))

	assert.Equal(t, cache.list.Tail(), cache.storage["third"])
	assert.Equal(t, cache.list.Head(), cache.storage["first"])

	// Exceeding the capacity
	cache.set("fourth", 10)

	assert.EqualValues(t, 3, cache.list.Size())
	assert.EqualValues(t, 3, len(cache.storage))

	assert.Equal(t, cache.list.Tail(), cache.storage["fourth"])
	assert.Equal(t, cache.list.Head(), cache.storage["second"])

	_, err = cache.get("first")
	assert.ErrorIs(t, ErrNotFound, err)

	_, err = cache.get("third")
	assert.NoError(t, err)

	assert.Equal(t, cache.list.Tail(), cache.storage["third"])
	assert.Equal(t, cache.list.Head(), cache.storage["second"])
}

func TestGetSetDataRace(t *testing.T) {
	cache := NewCache()

	cache.set("first", 1)

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()

			_, err := cache.get("first")
			assert.NoError(t, err)
			cache.set("first", 1)
		}()
	}

	wg.Wait()
}

func TestTestGetSetContext(t *testing.T) {
	cache := NewCache()

	ctx1, cancel := context.WithCancel(context.Background())
	cancel()

	err := cache.Set(ctx1, "1", 1)

	assert.ErrorIs(t, context.Canceled, err)

	_, err = cache.Get(ctx1, "1")

	assert.ErrorIs(t, context.Canceled, err)

	ctx2, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = cache.Set(ctx2, "1", 1)

	assert.NoError(t, err)

	v, err := cache.Get(ctx2, "1")

	assert.NoError(t, err)
	assert.Equal(t, 1, v)
}

func TestGetSetPanics(t *testing.T) {
	cache := NewCache()

	assert.Panics(t, func() {
		cache.Set(nil, "", 1)
	})

	assert.Panics(t, func() {
		cache.Get(nil, "")
	})
}
