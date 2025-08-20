package cache

var (
	_ Cache = (*cache)(nil)
)

type Cache interface {
	Set(key string, value any)
	Get(key string) (any, bool)
	Del(key string) any
	RPush(key string, data []any) int
	LPush(key string, data []any) int
	LRange(key string, start, end int) []any
	LLen(s string) int
	RPop(key string, count int) []any
	LPop(key string, count int) []any
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

func (c *cache) runJob() {
	for {
		select {}
	}
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

func (c *cache) LPush(key string, data []any) int {
	v, _ := c.listData[key]
	c.listData[key] = append(data, v...)
	return len(v) + len(data)
}

func (c *cache) LRange(key string, start, end int) []any {
	v, _ := c.listData[key]

	if len(v) == 0 {
		return []any{}
	}

	if start < 0 {
		if start < -len(v) {
			start = 0
		} else {
			start = len(v) + start
		}
	}

	if end < 0 {
		if end < -len(v) {
			return []any{}
		}

		end = len(v) + end
	}

	if start > end {
		return []any{}
	}

	end = end + 1
	if end > len(v) {
		end = len(v)
	}

	return v[start:end]
}

func (c *cache) LLen(s string) int {
	v, _ := c.listData[s]
	return len(v)
}

func (c *cache) RPop(key string, count int) []any {
	v, _ := c.listData[key]
	if len(v) == 0 {
		return nil
	}

	if count == 0 { // default
		count = 1
	}

	if count > len(v) {
		return v
	}

	endIdx := len(v) - 1 - count
	r := v[endIdx:]
	if count == len(v) {
		c.listData[key] = nil
		return r
	}

	c.listData[key] = v[:endIdx]
	return r
}

func (c *cache) LPop(key string, count int) []any {
	v, _ := c.listData[key]
	if len(v) == 0 {
		return nil
	}

	if count == 0 { // default
		count = 1
	}

	if count > len(v) {
		return v
	}

	r := v[0:count]
	if count == len(v) {
		c.listData[key] = nil
		return r
	}

	c.listData[key] = v[count:]
	return r
}
