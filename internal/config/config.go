package config

import (
	"errors"
	"strconv"
	"time"
)

// NonZeroUint64 a custom uint64 that cannot be zero.
type NonZeroUint64 uint64

// SetValue implements cleanenv.Setter interface.
func (u *NonZeroUint64) SetValue(s string) error {
	if s == "" {
		*u = 2048
		return nil
	}

	capacity, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}

	if capacity == 0 {
		return errors.New("CACHE_CAPACITY must be greater than 0")
	}

	*u = NonZeroUint64(capacity)
	return nil
}

// ToUint64 returns the value as uint64.
func (u NonZeroUint64) ToUint64() uint64 {
	return uint64(u)
}

// CacheConfig is the cache config struct.
type CacheConfig struct {
	CacheCapacity NonZeroUint64 `env:"CACHE_CAPACITY" env-default:"2048"`
}

type ServerConfig struct {
	Address      string        `env:"SERVER_ADDRESS" env-default:"127.0.0.1:2376"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" env-default:"1s"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" env-default:"1s"`
}
