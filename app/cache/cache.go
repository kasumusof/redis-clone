package cache

var (
	_ Cache = (*cache)(nil)
)

type Cache interface {
	Set(key string, value any)
	Get(key string) (any, bool)
	Del(key string) any
	RPush(key string, data []any) int
	LRange(key string, start, end int) []any
}
type cache struct {
	data     map[any]any
	listData map[any][]any
}

func New() Cache {
	c := &cache{
		data:     make(map[any]any),
		listData: make(map[any][]any),
	}

	go c.runJob()
	return c
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

func (c *cache) RPush(key string, data []any) int {
	v, _ := c.listData[key]
	c.listData[key] = append(v, data...)
	return len(v) + len(data)
}

func (c *cache) LRange(key string, start, end int) []any {
	end += 1
	v, _ := c.listData[key]

	if start < 0 {
		start = len(v) + start
	}
	if end < 0 {
		end = len(v) + end
	}

	if len(v) == 0 || start > end {
		return []any{}
	}

	if end > len(v) {
		end = len(v)
	}

	return v[start:end]
}

func (c *cache) runJob() {
	for {
		select {}
	}
}
