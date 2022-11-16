package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetValue(t *testing.T) {
	var u NonZeroUint64

	assert.NoError(t, u.SetValue(""))
	assert.EqualValues(t, 2048, u)

	assert.Error(t, u.SetValue("123.5"))

	assert.Equal(t, "CACHE_CAPACITY must be greater than 0", u.SetValue("0").Error())

	assert.NoError(t, u.SetValue("1000"))
	assert.EqualValues(t, 1000, u)
}

func TestToUint64(t *testing.T) {
	var u NonZeroUint64
	assert.NoError(t, u.SetValue(""))

	assert.Equal(t, uint64(2048), u.ToUint64())
}
