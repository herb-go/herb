package redispool

import (
	"encoding/json"
	"testing"
)

func TestRedis(t *testing.T) {
	redis := New()
	config := NewConfig()
	err := json.Unmarshal([]byte(testConfig), config)
	if err != nil {
		t.Fatal(err)
	}
	err = config.ApplyTo(redis)
	if err != nil {
		t.Fatal(err)
	}
	redis.Open()
	defer redis.Close()
	conn := redis.Get()
	defer conn.Close()

}

func TestConfig(t *testing.T) {
	redis := New()
	config := NewConfig()
	err := json.Unmarshal([]byte(testConfig), config)
	if err != nil {
		t.Fatal(err)
	}
	config.MaxIdle = 0
	config.MaxAlive = 0
	config.IdleTimeoutInSecond = 0
	err = config.ApplyTo(redis)
	if err != nil {
		t.Fatal(err)
	}
	redis.Open()
	defer redis.Close()
	if redis.Pool.MaxIdle != defaultMaxIdle {
		t.Fatal(redis.Pool.MaxIdle)
	}
	if redis.Pool.MaxActive != defaultMaxAlive {
		t.Fatal(redis.Pool.MaxActive)
	}
	if redis.Pool.IdleTimeout != defaultIdleTimeout {
		t.Fatal(redis.Pool.IdleTimeout)
	}
}
