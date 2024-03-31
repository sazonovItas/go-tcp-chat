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

	t.Run("check set and get", func(t *testing.T) {
		err := cache.Set(context.Background(), "test_set_and_get", "test", 0)
		if err != nil {
			t.Fatal(fmt.Errorf("%s: %w", "error to get value with a key", err))
		}

		value, err := cache.Get(context.Background(), "test_set_and_get")
		if err != nil {
			t.Fatal(fmt.Errorf("%s: %w", "error to get value with a key", err))
		}

		assert.Equal(t, "test", value, "should be equal get after set value")
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

	t.Run("check set struct", func(t *testing.T) {
		testData := testStruct{
			ID:      1,
			UUID:    uuid.New(),
			Product: "test",
			Price:   102.24,
		}

		err := cache.Set(context.Background(), "test_set_and_get_struct", testData, 0)
		if err != nil {
			t.Fatal(fmt.Errorf("%s: %w", "error to set value under a key", err))
		}

		value, err := cache.Get(context.Background(), "test_set_and_get_struct")
		if err != nil {
			t.Fatal(fmt.Errorf("%s: %w", "error to get value with a key", err))
		}

		assert.Equal(t, testData, value, "set and get value should be equal")
	})
}
