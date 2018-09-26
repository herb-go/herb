package cache

import "time"

//Field cache field component
type Field struct {
	Cache     Cacheable
	FieldName string
}

//Set set field value with given ttl.
//Return any error if raised.
func (f *Field) Set(v interface{}, ttl time.Duration) error {
	return f.Cache.Set(f.FieldName, v, ttl)
}

//Get value from field with given key.
//Return any error if raised.
func (f *Field) Get(key string, v interface{}) error {
	return f.Cache.Get(f.FieldName, v)
}

//SetBytesValue set bytes to field with given ttl.
//Return any error if raised.
func (f *Field) SetBytesValue(bytes []byte, ttl time.Duration) error {
	return f.Cache.SetBytesValue(f.FieldName, bytes, ttl)
}

//GetBytesValue get bytes value from field
//Return bytes value and any error if raised.
func (f *Field) GetBytesValue() ([]byte, error) {
	return f.Cache.GetBytesValue(f.FieldName)
}

//Del del field value.
//Return any error if raised.
func (f *Field) Del() error {
	return f.Cache.Del(f.FieldName)
}

//IncrCounter incr field counter with given increment and ttl
//Return new counter value and any error if raised.
func (f *Field) IncrCounter(increment int64, ttl time.Duration) (int64, error) {
	return f.Cache.IncrCounter(f.FieldName, increment, ttl)
}

//SetCounter set field counter with given value and ttl.
//Return any error if raised.
func (f *Field) SetCounter(v int64, ttl time.Duration) error {
	return f.Cache.SetCounter(f.FieldName, v, ttl)
}

//GetCounter get field counter.
//Return counter value and any error if raised.
func (f *Field) GetCounter() (int64, error) {
	return f.Cache.GetCounter(f.FieldName)
}

//DelCounter delete field counter.
//Return any error if raised.
func (f *Field) DelCounter() error {
	return f.Cache.DelCounter(f.FieldName)
}
