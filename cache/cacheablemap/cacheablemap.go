package cacheablemap

import (
	"reflect"

	"github.com/herb-go/herb/cache"
)

//Map cached map interface
type Map interface {
	//NewMapElement method create new Element with given key.
	//Return any error if raised.
	NewMapElement(string) error
	//LoadMapElements method load element to map by give key list.
	//Return any error if raised.
	LoadMapElements(keys ...string) error
	//Map return cachable map
	Map() interface{}
}

func unmarshalMapElement(cm interface{}, creator func(string) error, key string, data []byte) (err error) {
	err = creator(key)
	if err != nil {
		return err
	}
	var mapvalue = reflect.Indirect(reflect.ValueOf(cm))
	var v = mapvalue.MapIndex(reflect.ValueOf(key))
	var vp = reflect.New(v.Type())
	vp.Elem().Set(v)
	err = cache.UnmarshalMsgpack(data, vp.Interface())
	if err != nil {
		return err
	}
	mapvalue.SetMapIndex(reflect.ValueOf(key), vp.Elem())
	return nil
}

//Load load data by given key list into cache map.
//Param cm target data map pointer.
//Param c map cache
//Param loader func by which data load.
//Param creator map element creator.
//Param keys key list to load
//Return any error if raised.
func Load(cm interface{}, c cache.Cacheable, loader func(keys ...string) error, creator func(string) error, keys ...string) error {
	var keysmap = make(map[string]bool, len(keys))
	var mapvalue = reflect.Indirect(reflect.ValueOf(cm))
	var filteredKeys = make([]string, len(keys))
	var filteredKeysLength = 0
	if len(keys) == 0 {
		return nil
	}
	for k := range keys {
		if keysmap[keys[k]] == true {
			continue
		}
		keysmap[keys[k]] = true
		if !mapvalue.MapIndex(reflect.ValueOf(keys[k])).IsValid() {

			filteredKeys[filteredKeysLength] = keys[k]
			filteredKeysLength++
		}
	}
	filteredKeys = filteredKeys[:filteredKeysLength]
	if filteredKeysLength == 0 {
		return nil
	}
	lockKey := filteredKeys[0]
	_, err := c.Wait(lockKey)
	if err != nil {
		return err
	}
	results, err := c.MGetBytesValue(filteredKeys...)
	if err != nil {
		return err
	}
	var uncachedKeys = make([]string, len(filteredKeys))
	var uncachedKeysLength = 0
	for i := range filteredKeys {
		k := filteredKeys[i]
		if results[k] == nil {
			err = creator(k)
			if err != nil {
				return err
			}
			uncachedKeys[uncachedKeysLength] = k
			uncachedKeysLength++
		} else {
			err = unmarshalMapElement(cm, creator, k, results[k])
			if err != nil {
				return err
			}
		}
	}
	uncachedKeys = uncachedKeys[:uncachedKeysLength]
	if uncachedKeysLength == 0 {
		return nil
	}
	unlocker, err := c.Lock(lockKey)
	if err != nil {
		return err
	}
	defer unlocker()
	err = loader(uncachedKeys...)
	if err != nil {
		return err
	}
	var data = make(map[string][]byte, len(uncachedKeys))
	for k := range uncachedKeys {
		v := mapvalue.MapIndex(reflect.ValueOf(uncachedKeys[k])).Interface()
		data[uncachedKeys[k]], err = cache.MarshalMsgpack(v)
		if err != nil {
			return err
		}
	}
	return c.MSetBytesValue(data, 0)
}

//LoadCachedMap load data by given key list into cachedmap interface.
//Param cm target data map pointer.
//Param c map cache
//Param keys key list to load
//Return any error if raised.
func LoadCachedMap(cm Map, c cache.Cacheable, keys ...string) error {
	return Load(cm.Map(), c, cm.LoadMapElements, cm.NewMapElement, keys...)
}
