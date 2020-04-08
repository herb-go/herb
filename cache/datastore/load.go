package datastore

import (
	"sort"

	"github.com/herb-go/herb/cache"
)

func unmarshalMapElement(s Store, creator func() interface{}, key string, data []byte, c cache.Cacheable) (err error) {
	v := creator()
	if len(data) != 0 {
		err = c.Util().Unmarshal(data, v)
		if err != nil {
			return err
		}
	} else {
		v = nil
	}
	s.Store(key, v)
	return nil
}

//Load load data by given key list into data store.
//Param s target data store.
//Param c map cache
//Param loader func by which data load.
//Param creator map element creator.
//Param keys key list to load
//Return any error if raised.
func Load(s Store, c cache.Cacheable, loader func(...string) (map[string]interface{}, error), creator func() interface{}, keys ...string) error {
	var keysmap = make(map[string]bool, len(keys))
	var filteredKeys = make([]string, 0, len(keys))
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

			filteredKeys = append(filteredKeys, keys[k])
		}
	}
	if len(filteredKeys) == 0 {
		return nil
	}
	var results map[string][]byte
	if c != nil {

		results, err = c.MGetBytesValue(filteredKeys...)
		if err != nil {
			return err
		}
	} else {
		results = map[string][]byte{}
	}
	var uncachedKeys = make([]string, 0, len(filteredKeys))
	for i := range filteredKeys {
		k := filteredKeys[i]
		if _, ok := results[k]; !ok {
			uncachedKeys = append(uncachedKeys, k)
		} else {
			err = unmarshalMapElement(s, creator, k, results[k], c)
			if err != nil {
				return err
			}
		}
	}
	if len(uncachedKeys) == 0 {
		return nil
	}
	var unreloadkeys = make([]string, 0, len(uncachedKeys))

	sort.Strings(uncachedKeys)
	if c != nil {
		for k := range uncachedKeys {
			key := c.FinalKey(uncachedKeys[k])

			locker, ok := c.Util().Locker(key)
			if ok {
				locker.RLock()
				defer locker.RUnlock()
				bs, err := c.GetBytesValue(uncachedKeys[k])
				if err == nil {
					err = unmarshalMapElement(s, creator, uncachedKeys[k], bs, c)
					if err != nil {
						return err
					}
					continue
				}
				if err != cache.ErrNotFound {
					return err
				}
			} else {
				locker.Lock()
				defer locker.Unlock()
			}
			unreloadkeys = append(unreloadkeys, uncachedKeys[k])

		}

	} else {
		unreloadkeys = uncachedKeys
	}

	if len(unreloadkeys) == 0 {
		return nil
	}
	loaded, err := loader(unreloadkeys...)
	if err != nil {
		return err
	}
	var data = make(map[string][]byte, len(loaded))
	for k := range loaded {
		v := loaded[k]
		s.Store(k, v)
		if c != nil {
			data[k], err = c.Util().Marshal(v)
			if err != nil {
				return err
			}
		}
	}
	for k := range unreloadkeys {
		if _, ok := data[unreloadkeys[k]]; ok == false {
			data[unreloadkeys[k]] = []byte{}
		}
	}
	if c == nil {
		return nil
	}
	return c.MSetBytesValue(data, 0)
}
