package cache_test

import (
	"testing"

	"github.com/herb-go/herb/cache"
	_ "github.com/herb-go/herb/cache/drivers/syncmapcache"
)

func newFieldTestCache(ttl int64) *cache.Field {
	config := &cache.ConfigMap{}
	config.Set("Size", 10000000)
	c := cache.New()
	oc := &cache.OptionConfigMap{
		Driver:    "syncmapcache",
		TTL:       int64(ttl),
		Config:    nil,
		Marshaler: "json",
	}
	err := c.Init(oc)
	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	return c.Field("testnode")
}

func TestFieldUpdate(t *testing.T) {
	var err error
	defaultTTL := int64(1)
	c := newFieldTestCache(defaultTTL)
	testDataModel := "test"
	var resultDataModel string
	testDataBytes := []byte("testbytes")
	err = c.Set(testDataModel, cache.TTLForever)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(&resultDataModel)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Update(testDataModel, cache.TTLForever)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(&resultDataModel)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Del()
	if err != nil {
		t.Fatal(err)
	}
	err = c.Update(testDataModel, cache.TTLForever)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(&resultDataModel)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}

	err = c.SetBytesValue(testDataBytes, cache.TTLForever)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue()
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetBytesValue(testDataBytes, cache.TTLForever)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue()
	if err != nil {
		t.Fatal(err)
	}
	err = c.Del()
	if err != nil {
		t.Fatal(err)
	}
	err = c.UpdateBytesValue(testDataBytes, cache.TTLForever)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue()
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
}

func TestFieldCounter(t *testing.T) {
	defaultTTL := int64(1)
	testInitVal := int64(1)
	testIncremeant := int64(2)
	testTargetResultInt := int64(3)
	var resultDataInt int64
	c := newFieldTestCache(defaultTTL)
	err := c.SetCounter(testInitVal, cache.DefualtTTL)
	if err != nil {
		t.Fatal(err)
	}
	resultDataInt, err = c.GetCounter()
	if err != nil {
		t.Fatal(err)
	}
	if resultDataInt != testInitVal {
		t.Errorf("GetCounter error %d ", resultDataInt)
	}
	resultDataInt, err = c.IncrCounter(testIncremeant, cache.DefualtTTL)
	if err != nil {
		t.Fatal(err)
	}
	if resultDataInt != testTargetResultInt {
		t.Errorf("IncrCounter error %d ", resultDataInt)
	}
	resultDataInt, err = c.GetCounter()
	if err != nil {
		t.Fatal(err)
	}
	if resultDataInt != testTargetResultInt {
		t.Errorf("GetCounter error %d ", resultDataInt)
	}
	err = c.DelCounter()
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetCounter()
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
}
