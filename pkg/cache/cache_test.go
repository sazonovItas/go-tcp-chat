package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

const (
	testAddr     = ":6379"
	testPassword = ""
	testDB       = 0
)

func TestNewCache(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     testAddr,
		Password: testPassword,
		DB:       testDB,
	})

	if client.Ping(context.Background()).Err() != nil {
		return
	}

	cache := NewCache[string](&CacheOpts{
		Client: client,

		DefaultExpiration: time.Second * 15,
	})

	t.Run("check set data under a key", func(t *testing.T) {
		err := cache.Set(context.Background(), "test_data", "test", 0)
		if err != nil {
			t.Fatal(fmt.Errorf("%s: %w", "error to get value with a key", err))
		}
	})

	t.Run("check exists", func(t *testing.T) {
		ok := cache.Exists(context.Background(), "test_data")
		assert.Equal(t, true, ok, "key should exists")
	})

	t.Run("check get data with a key", func(t *testing.T) {
		value, err := cache.Get(context.Background(), "test_data")
		if err != nil {
			t.Fatal(fmt.Errorf("%s: %w", "error to get value with a key", err))
		}

		assert.Equal(t, "test", value, "should be equal get after set value")
	})

	t.Run("check delete data under a key", func(t *testing.T) {
		err := cache.Delete(context.Background(), "test_data")
		assert.Equal(t, nil, err, "should not be error delete data under a key from cache")
	})

	t.Run("check get data with non-existent key", func(t *testing.T) {
		_, err := cache.Get(context.Background(), "test_data")
		if assert.Error(t, err, "should be error to retrive non-existent key") {
			assert.Equal(t, redis.Nil, err, "should nil error to retrive non-existent data")
		}
	})
}

func TestSetGetStructures(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     testAddr,
		Password: testPassword,
		DB:       testDB,
	})

	if client.Ping(context.Background()).Err() != nil {
		return
	}

	type testStruct struct {
		ID        int64     `json:"id"`
		UUID      uuid.UUID `json:"uuid"`
		Product   string    `json:"product"`
		Price     float64   `json:"price"`
		CreatedAt time.Time `json:"created_at"`
	}

	cache := NewCache[testStruct](&CacheOpts{
		Client: client,

		DefaultExpiration: time.Second * 15,
	})

	testData := testStruct{
		ID:      1,
		UUID:    uuid.New(),
		Product: "test",
		Price:   102.24,
	}

	t.Run("check set data under a key", func(t *testing.T) {
		err := cache.Set(context.Background(), "test_struct", testData, 0)
		if err != nil {
			t.Fatal(fmt.Errorf("%s: %w", "error to set value under a key", err))
		}
	})

	t.Run("check get value", func(t *testing.T) {
		value, err := cache.Get(context.Background(), "test_struct")
		if err != nil {
			t.Fatal(fmt.Errorf("%s: %w", "error to get value with a key", err))
		}

		assert.Equal(t, testData, value, "set and get value should be equal")
	})

	t.Run("check get expired data", func(t *testing.T) {
		err := cache.Set(context.Background(), "test_expired", testStruct{}, time.Millisecond)
		if err != nil {
			t.Fatal(fmt.Errorf("%s: %w", "error to set value under a key", err))
		}

		time.Sleep(time.Millisecond * 20)

		_, err = cache.Get(context.Background(), "test_expired")
		if assert.Error(t, err, "should be error to retrive expired key") {
			assert.Equal(t, redis.Nil, err, "should nil error to retrive expired data")
		}
	})
}
