package cache

type Proxy struct {
	Cacheable
}

func NewProxy(c Cacheable) *Proxy {
	return &Proxy{
		Cacheable: c,
	}
}

func ProxyWithPrefix(c Cacheable, prefix string) *Proxy {
	return NewProxy(NewCollection(c, prefix, DefaultTTL))
}
