package cache

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

func Set(key string, value any) {
	defaultCache.Set(key, value)
}

func Get(key string) (any, bool) {
	return defaultCache.Get(key)
}
