package cachestore

import "github.com/herb-go/herb/cache"

func unmarshalMapElement(s Store, creator func() interface{}, key string, data []byte, c cache.Cacheable) (err error) {
	v := creator()
	if len(data) != 0 {
		err = c.Unmarshal(data, v)
		if err != nil {
			return err
		}
	} else {
		v = nil
	}
	s.Store(key, v)
	return nil
}

//Load load data by given key list into cache store.
//Param s target data map pointer.
//Param c map cache
//Param loader func by which data load.
//Param creator map element creator.
//Param keys key list to load
//Return any error if raised.
func Load(s Store, c cache.Cacheable, loader func(...string) (map[string]interface{}, error), creator func() interface{}, keys ...string) error {
	var keysmap = make(map[string]bool, len(keys))
	var filteredKeys = make([]string, len(keys))
	var filteredKeysLength = 0
	var err error
	if len(keys) == 0 {
		return nil
	}
	for k := range keys {
		if keysmap[keys[k]] == true {
			continue
		}
		keysmap[keys[k]] = true
		if _, ok := s.Load(keys[k]); !ok {

			filteredKeys[filteredKeysLength] = keys[k]
			filteredKeysLength++
		}
	}
	filteredKeys = filteredKeys[:filteredKeysLength]
	if filteredKeysLength == 0 {
		return nil
	}
	var results map[string][]byte
	lockers := map[string]*cache.Locker{}
	if c != nil {
		for k := range filteredKeys {
			key, err := c.FinalKey(filteredKeys[k])
			if err != nil {
				return err
			}
			locker, err := c.Locker(key)
			if err != nil {
				return err
			}
			locker.RLock()
			lockers[filteredKeys[k]] = locker
			defer locker.Unlock()
		}

		results, err = c.MGetBytesValue(filteredKeys...)
		if err != nil {
			return err
		}
	} else {
		results = map[string][]byte{}
	}
	var uncachedKeys = make([]string, len(filteredKeys))
	var uncachedKeysLength = 0
	for i := range filteredKeys {
		k := filteredKeys[i]
		if _, ok := results[k]; !ok {
			uncachedKeys[uncachedKeysLength] = k
			uncachedKeysLength++
		} else {
			err = unmarshalMapElement(s, creator, k, results[k], c)
			if err != nil {
				return err
			}
		}
	}
	uncachedKeys = uncachedKeys[:uncachedKeysLength]
	if uncachedKeysLength == 0 {
		return nil
	}
	if c != nil {
		for k := range uncachedKeys {
			locker := lockers[filteredKeys[k]]
			locker.Lock()
		}

	}
	loaded, err := loader(uncachedKeys...)
	if err != nil {
		return err
	}
	var data = make(map[string][]byte, len(loaded))
	for k := range loaded {
		v := loaded[k]
		s.Store(k, v)
		if c != nil {
			data[k], err = c.Marshal(v)
			if err != nil {
				return err
			}
		}
	}
	for k := range uncachedKeys {
		if _, ok := data[uncachedKeys[k]]; ok == false {
			data[uncachedKeys[k]] = []byte{}
		}
	}
	if c == nil {
		return nil
	}
	return c.MSetBytesValue(data, 0)
}
