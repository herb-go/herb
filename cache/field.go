package cache

import "time"

type Field struct {
	Cache     *Cache
	FieldName string
}

func (f *Field) Set(v interface{}, ttl time.Duration) error {
	return f.Cache.Set(f.FieldName, v, ttl)
}

func (f *Field) Update(v interface{}, ttl time.Duration) error {
	return f.Cache.Update(f.FieldName, v, ttl)
}

func (f *Field) Get(key string, v interface{}) error {
	return f.Cache.Get(f.FieldName, &v)
}

func (f *Field) SetBytesValue(bytes []byte, ttl time.Duration) error {
	return f.Cache.SetBytesValue(f.FieldName, bytes, ttl)
}

func (f *Field) UpdateBytesValue(bytes []byte, ttl time.Duration) error {
	return f.Cache.UpdateBytesValue(f.FieldName, bytes, ttl)
}

func (f *Field) GetBytesValue() ([]byte, error) {
	return f.Cache.GetBytesValue(f.FieldName)
}

func (f *Field) Del() error {
	return f.Cache.Del(f.FieldName)
}

func (f *Field) IncrCounter(increment int64, ttl time.Duration) (int64, error) {
	return f.Cache.IncrCounter(f.FieldName, increment, ttl)
}
func (f *Field) SetCounter(v int64, ttl time.Duration) error {
	return f.Cache.SetCounter(f.FieldName, v, ttl)
}

func (f *Field) GetCounter() (int64, error) {
	return f.Cache.GetCounter(f.FieldName)
}
func (f *Field) DelCounter() error {
	return f.Cache.DelCounter(f.FieldName)
}
