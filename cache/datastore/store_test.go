package datastore

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/herb-go/herb/cache"
	_ "github.com/herb-go/herb/cache/drivers/syncmapcache"
)

func newTestCache(ttl int64) *cache.Cache {
	config := &cache.ConfigJSON{}
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
	return c

}

type testmodel struct {
	Keyword string
	Content int
}

const valueKey = "valueKey"
const valueKeyAadditional = "valueKeyAadditional"
const valueKeyChanged = "valueKeyChanged"
const valueKeyNotexists = "valueKeyNotExists"
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
			raw, ok := rawData[v]
			if ok == false {
				continue
			}
			result[v] = &testmodel{
				Keyword: creatorKeyword,
				Content: raw,
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
	if v == nil {
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
	err = Load(tm, c, loader(), creator(), valueKey, valueKey, valueKeyAadditional, valueKeyNotexists)
	if err != nil {
		t.Fatal(err)
	}
	if val := load(tm, valueKeyNotexists); val != nil {
		t.Fatal(val)
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
	rawData[valueKeyNotexists] = changedValue
	rawData[valueKey] = changedValue
	rawData[valueKeyAadditional] = changedValue
	rawData[valueKeyChanged] = changedValue

	err = Load(tm, c, loader(), creator(), valueKeyAadditional, valueKeyChanged, valueKeyNotexists)
	if err != nil {
		t.Fatal(err)
	}
	if val := load(tm, valueKeyNotexists); val != nil {
		t.Error(val)
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

func TestNilcache(t *testing.T) {
	rawData = map[string]int{
		valueKey:            startValue,
		valueKeyAadditional: startValue,
		valueKeyChanged:     startValue,
	}
	var c cache.Cacheable
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

func TestLoader(t *testing.T) {
	rawData = map[string]int{
		valueKey:            startValue,
		valueKeyAadditional: startValue,
		valueKeyChanged:     startValue,
	}
	c := newTestCache(-1)
	var err error
	var datasource = NewDataSource()
	datasource.SourceLoader = loader()
	datasource.Creator = creator()
	var Loader = datasource.NewMapStoreLoader(c)
	tm := Loader.Store
	err = Loader.Load(valueKey, valueKeyAadditional)
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
	err = Loader.Load(valueKeyAadditional, valueKeyChanged)
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
	val := Loader.Store.LoadInterface(valueKeyAadditional)
	if val == nil {
		t.Fatal(val)
	}
	val = Loader.Store.LoadInterface(valueKeyChanged)
	if val == nil {
		t.Fatal(val)
	}
	val = Loader.Store.LoadInterface(valueKeyNotexists)
	if val != nil {
		t.Fatal(val)
	}
	err = Loader.Del(valueKeyAadditional)
	val = Loader.Store.LoadInterface(valueKeyAadditional)
	if val != nil {
		t.Fatal(val)
	}
	err = Loader.Flush()
	val = Loader.Store.LoadInterface(valueKeyChanged)
	if val != nil {
		t.Fatal(val)
	}
}

func TestSyncLoader(t *testing.T) {
	rawData = map[string]int{
		valueKey:            startValue,
		valueKeyAadditional: startValue,
		valueKeyChanged:     startValue,
	}
	c := newTestCache(-1)
	var err error
	var datasource = &DataSource{
		SourceLoader: loader(),
		Creator:      creator(),
	}
	var Loader = datasource.NewSyncMapStoreLoader(c)
	tm := Loader.Store
	err = Loader.Load(valueKey, valueKeyAadditional)
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
	err = Loader.Load(valueKeyAadditional, valueKeyChanged)
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
	val := Loader.Store.LoadInterface(valueKeyAadditional)
	if val == nil {
		t.Fatal(val)
	}
	val = Loader.Store.LoadInterface(valueKeyChanged)
	if val == nil {
		t.Fatal(val)
	}
	val = Loader.Store.LoadInterface(valueKeyNotexists)
	if val != nil {
		t.Fatal(val)
	}
	err = Loader.Del(valueKeyAadditional)
	val = Loader.Store.LoadInterface(valueKeyAadditional)
	if val != nil {
		t.Fatal(val)
	}
	err = Loader.Flush()
	val = Loader.Store.LoadInterface(valueKeyChanged)
	if val != nil {
		t.Fatal(val)
	}
}

func TestConcurrent(t *testing.T) {
	rawdata := sync.Map{}
	c := newTestCache(-1)
	loader := func(keys ...string) (map[string]interface{}, error) {
		result := map[string]interface{}{}
		time.Sleep(300 * time.Microsecond)
		for _, v := range keys {
			var data int
			val, _ := rawdata.Load(v)
			if val == nil {
				data = 0
			} else {
				data = val.(int)
			}
			data++
			rawdata.Store(v, data)
			result[v] = &data
		}
		return result, nil
	}
	creator := func() interface{} {
		v := int(0)
		return &v
	}
	var tm = NewSyncMapStore()
	var tm2 = NewSyncMapStore()
	var tm3 = NewSyncMapStore()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Load(tm, c, loader, creator, valueKey, valueKeyAadditional)
		if err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(100 * time.Microsecond)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Load(tm2, c, loader, creator, valueKey, valueKeyChanged)
		if err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(100 * time.Microsecond)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Load(tm3, c, loader, creator, valueKey, valueKeyNotexists)
		if err != nil {
			t.Fatal(err)
		}
	}()

	wg.Wait()
	if *(tm.LoadInterface(valueKey).(*int)) != 1 {
		t.Fatal(tm)
	}
	if *(tm.LoadInterface(valueKeyAadditional).(*int)) != 1 {
		t.Fatal(tm)
	}
	if *(tm2.LoadInterface(valueKey).(*int)) != 1 {
		t.Fatal(tm2)
	}
	if *(tm2.LoadInterface(valueKeyChanged).(*int)) != 1 {
		t.Fatal(tm2)
	}
	if *(tm3.LoadInterface(valueKey).(*int)) != 1 {
		t.Fatal(tm3)
	}
	if *(tm3.LoadInterface(valueKeyNotexists).(*int)) != 1 {
		t.Fatal(tm3)
	}
	locker, _ := c.Util().Locker("test")
	locker.Lock()
	locker.Unlock()
	locker.Map.Range(func(key interface{}, value interface{}) bool {
		//check unlocked locker
		t.Fatal(key, value)
		return true
	})
}

func TestMutliConcurrent(t *testing.T) {
	rawdata := sync.Map{}
	c := newTestCache(-1)
	loader := func(keys ...string) (map[string]interface{}, error) {
		result := map[string]interface{}{}
		time.Sleep(10 * time.Millisecond)
		for _, v := range keys {
			var data int
			val, _ := rawdata.Load(v)
			if val == nil {
				data = 0
			} else {
				data = val.(int)
			}
			data++
			rawdata.Store(v, data)
			result[v] = &data
		}
		return result, nil
	}
	creator := func() interface{} {
		v := int(0)
		return &v
	}
	wg := &sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			keys := []string{}
			length := rand.Intn(10)
			for n := 0; n < length; n++ {
				keys = append(keys, strconv.Itoa(rand.Intn(100)))
			}
			tm := NewSyncMapStore()
			err := Load(tm, c, loader, creator, keys...)
			if err != nil {
				t.Fatal(err)
			}
		}()
	}
	wg.Wait()
	locker, _ := c.Util().Locker("test")
	locker.Lock()
	locker.Unlock()
	locker.Map.Range(func(key interface{}, value interface{}) bool {
		//check unlocked locker
		t.Fatal(key, value)
		return true
	})

}
