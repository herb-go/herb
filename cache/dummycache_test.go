package cache_test

import (
	"testing"
	"time"

	"github.com/herb-go/herb/cache"
)

func TestDummyCache(t *testing.T) {
	c := cache.New()
	testKey := "key"
	testData := []string{"123"}
	testBytes := []byte("123")
	testTTL := time.Hour
	testIncrement := int64(1)
	testIntValue := int64(2)
	var model string
	oc := cache.NewOptionConfig()
	oc.Driver = "dummycache"
	oc.TTL = int64(testTTL)
	oc.Config = nil
	oc.Marshaler = "json"
	err := c.Init(oc)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Set(testKey, testData, testTTL)
	if err != nil {
		t.Errorf("Set error %s", err)
	}
	err = c.Update(testKey, testData, testTTL)
	if err != nil {
		t.Errorf("Update error %s", err)
	}
	err = c.Get(testKey, &model)
	if err != cache.ErrNotFound {
		t.Errorf("Get error %s", err)
	}
	err = c.SetBytesValue(testKey, testBytes, testTTL)
	if err != nil {
		t.Errorf("SetBytesValue error %s", err)
	}
	err = c.UpdateBytesValue(testKey, testBytes, testTTL)
	if err != nil {
		t.Errorf("UpdateBytesValue error %s", err)
	}
	_, err = c.GetBytesValue(testKey)
	if err != cache.ErrNotFound {
		t.Errorf("GetBytesValue error %s", err)
	}
	r, err := c.MGetBytesValue("test", "test2")
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 0 {
		t.Fatal(err)
	}
	err = c.MSetBytesValue(map[string][]byte{}, testTTL)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Del(testKey)
	if err != nil {
		t.Errorf("Del error %s", err)
	}
	_, err = c.IncrCounter(testKey, testIncrement, testTTL)
	if err != nil {
		t.Errorf("IncrCounter error %s", err)
	}
	err = c.SetCounter(testKey, testIntValue, testTTL)
	if err != nil {
		t.Errorf("SetCounter error %s", err)
	}
	_, err = c.GetCounter(testKey)
	if err != cache.ErrNotFound {
		t.Errorf("GetCounter error %s", err)
	}
	err = c.DelCounter(testKey)
	if err != nil {
		t.Errorf("DelCounter error %s", err)
	}
	err = c.Flush()
	if err != nil {
		t.Errorf("Flush error %s", err)
	}
	err = c.Expire(testKey, testTTL)
	if err != nil {
		t.Errorf("Flush error %s", err)
	}
	err = c.ExpireCounter(testKey, testTTL)
	if err != nil {
		t.Errorf("Flush error %s", err)
	}
	c.SetGCErrHandler(nil)
	err = c.Close()
	if err != nil {
		t.Errorf("Close error %s", err)
	}

}
