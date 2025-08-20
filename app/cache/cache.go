package cache

import "time"

var (
	_            Cache = (*cache)(nil)
	defaultCache Cache
)

func init() {
	defaultCache = New()
}

type Cache interface {
	Set(key string, value any)
	Get(key string) (any, bool)
	Del(key string) any
}
type cache struct {
	data map[any]any
}

func (c *cache) Set(key string, value any) {
	c.data[key] = value
}

func (c *cache) Get(key string) (any, bool) {
	val, ok := c.data[key]
	return val, ok
}

func (c *cache) Del(key string) any {
	old := c.data[key]
	delete(c.data, key)
	return old
}

func New() Cache {
	c := &cache{
		data: make(map[any]any),
	}

	go c.runJob()
	return c
}

func (c *cache) runJob() {
	for {
		select {}
	}
}

func Set(key string, value any, expiration int) {
	defaultCache.Set(key, value)
	if expiration > 0 {
		go func() {
			<-time.After(time.Duration(expiration) * time.Millisecond)
			Del(key)
		}()
	}
}

func Get(key string) (any, bool) {
	return defaultCache.Get(key)
}

func Del(key string) any {
	return defaultCache.Del(key)
}
