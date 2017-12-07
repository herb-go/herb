package cache

import (
	"time"
)

type Node struct {
	Cache  Cacheable
	Prefix string
}

func NewNode(c Cacheable, prefix string) *Node {
	return &Node{
		Cache:  c,
		Prefix: prefix,
	}
}
func (n *Node) GetCacheKey(key string) (string, error) {
	return n.Prefix + keyPrefix + key, nil
}
func (n *Node) MustGetCacheKey(key string) string {
	k, _ := n.GetCacheKey(key)
	return k
}

func (n *Node) Set(key string, v interface{}, ttl time.Duration) error {
	k := n.MustGetCacheKey(key)
	return n.Cache.Set(k, v, ttl)
}
func (n *Node) Get(key string, v interface{}) error {
	k := n.MustGetCacheKey(key)
	return n.Cache.Get(k, &v)
}
func (n *Node) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	k := n.MustGetCacheKey(key)
	return n.Cache.SetBytesValue(k, bytes, ttl)
}
func (n *Node) GetBytesValue(key string) ([]byte, error) {
	k := n.MustGetCacheKey(key)
	return n.Cache.GetBytesValue(k)
}
func (n *Node) Del(key string) error {
	k := n.MustGetCacheKey(key)
	return n.Cache.Del(k)
}
func (n *Node) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	k := n.MustGetCacheKey(key)
	return n.Cache.IncrCounter(k, increment, ttl)
}
func (n *Node) SetCounter(key string, v int64, ttl time.Duration) error {
	k := n.MustGetCacheKey(key)
	return n.Cache.SetCounter(k, v, ttl)
}
func (n *Node) GetCounter(key string) (int64, error) {
	k := n.MustGetCacheKey(key)
	return n.Cache.GetCounter(k)
}
func (n *Node) Load(key string, v interface{}, ttl time.Duration, closure func(v interface{}) error) error {
	k := n.MustGetCacheKey(key)
	return n.Cache.Load(k, v, ttl, closure)
}
func (n *Node) Flush() error {
	return ErrFeatureNotSupported
}
func (n *Node) DefualtTTL() time.Duration {
	return n.Cache.DefualtTTL()
}
func (n *Node) DelCounter(key string) error {
	k, err := n.GetCacheKey(key)
	if err != nil {
		return err
	}
	return n.Cache.DelCounter(k)
}

func (n *Node) Collection(prefix string) *Collection {
	return NewCollection(n, prefix, n.Cache.DefualtTTL())
}
func (n *Node) Node(prefix string) *Node {
	return NewNode(n, prefix)
}
func (n *Node) Field(fieldname string) *Field {
	return &Field{
		Cache:     n,
		FieldName: fieldname,
	}
}