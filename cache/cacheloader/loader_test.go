package cacheloader

import (
	"testing"

	"github.com/herb-go/herb/cache"
	_ "github.com/herb-go/herb/cache/drivers/freecache"
)

func newTestCache(ttl int64) *cache.Cache {
	config := &cache.ConfigJSON{}
	config.Set("Size", 10000000)
	c := cache.New()
	oc := &cache.OptionConfigMap{
		Driver:    "freecache",
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
	return c

}

type testmodel struct {
	Keyword string
	Content int
}

const valueKey = "valueKey"
const valueKeyAadditional = "valueKeyAadditional"
const valueKeyChanged = "valueKeyChanged"
const wrongDataKey = "wrongdata"

var WrongData = []byte("wrongdata")

const startValue = 1
const changedValue = 2
const mapCreatorKeyword = "mapCreatorKeyword"

const creatorKeyword = "creatorKeyword"

var rawData map[string]int

type testmodelmap map[string]*testmodel

func creator() func() interface{} {
	return func() interface{} {
		return &testmodel{
			Keyword: creatorKeyword,
			Content: 0,
		}
	}
}

func loader() func(...string) (map[string]interface{}, error) {
	return func(keys ...string) (map[string]interface{}, error) {
		var result = map[string]interface{}{}
		for _, v := range keys {
			result[v] = &testmodel{
				Keyword: creatorKeyword,
				Content: rawData[v],
			}
		}
		return result, nil
	}
}
func load(s Store, key string) *testmodel {
	v, ok := s.Load(key)
	if ok == false {
		return nil
	}
	return v.(*testmodel)
}
func TestMapLoad(t *testing.T) {
	rawData = map[string]int{
		valueKey:            startValue,
		valueKeyAadditional: startValue,
		valueKeyChanged:     startValue,
	}
	c := newTestCache(-1)
	var err error
	var tm = NewMapStore()
	err = Load(tm, c, loader(), creator(), valueKey, valueKeyAadditional)
	if err != nil {
		t.Fatal(err)
	}
	if val := load(tm, valueKey).Content; val != startValue {
		t.Error(val)
	}
	if val := load(tm, valueKey).Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := load(tm, valueKeyAadditional).Content; val != startValue {
		t.Error(val)
	}
	if val := load(tm, valueKeyAadditional).Keyword; val != creatorKeyword {
		t.Error(val)
	}
	rawData[valueKey] = changedValue
	rawData[valueKeyAadditional] = changedValue
	rawData[valueKeyChanged] = changedValue
	err = Load(tm, c, loader(), creator(), valueKeyAadditional, valueKeyChanged)
	if err != nil {
		t.Fatal(err)
	}
	if val := load(tm, valueKey).Content; val != startValue {
		t.Error(val)
	}
	if val := load(tm, valueKey).Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := load(tm, valueKeyAadditional).Content; val != startValue {
		t.Error(val)
	}
	if val := load(tm, valueKeyAadditional).Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := load(tm, valueKeyChanged).Content; val != changedValue {
		t.Error(val)
	}
	if val := load(tm, valueKeyChanged).Keyword; val != creatorKeyword {
		t.Error(val)
	}
	var tm2 = NewMapStore()
	err = Load(tm2, c, loader(), creator(), valueKeyAadditional, valueKeyChanged)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSyncMapLoad(t *testing.T) {
	rawData = map[string]int{
		valueKey:            startValue,
		valueKeyAadditional: startValue,
		valueKeyChanged:     startValue,
	}
	c := newTestCache(-1)
	var err error
	var tm = NewSyncMapStore()
	err = Load(tm, c, loader(), creator(), valueKey, valueKeyAadditional)
	if err != nil {
		t.Fatal(err)
	}
	if val := load(tm, valueKey).Content; val != startValue {
		t.Error(val)
	}
	if val := load(tm, valueKey).Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := load(tm, valueKeyAadditional).Content; val != startValue {
		t.Error(val)
	}
	if val := load(tm, valueKeyAadditional).Keyword; val != creatorKeyword {
		t.Error(val)
	}
	rawData[valueKey] = changedValue
	rawData[valueKeyAadditional] = changedValue
	rawData[valueKeyChanged] = changedValue
	err = Load(tm, c, loader(), creator(), valueKeyAadditional, valueKeyChanged)
	if err != nil {
		t.Fatal(err)
	}
	if val := load(tm, valueKey).Content; val != startValue {
		t.Error(val)
	}
	if val := load(tm, valueKey).Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := load(tm, valueKeyAadditional).Content; val != startValue {
		t.Error(val)
	}
	if val := load(tm, valueKeyAadditional).Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := load(tm, valueKeyChanged).Content; val != changedValue {
		t.Error(val)
	}
	if val := load(tm, valueKeyChanged).Keyword; val != creatorKeyword {
		t.Error(val)
	}
	var tm2 = NewMapStore()
	err = Load(tm2, c, loader(), creator(), valueKeyAadditional, valueKeyChanged)
	if err != nil {
		t.Fatal(err)
	}
}
