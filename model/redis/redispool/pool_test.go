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
	redis.Close()
}
