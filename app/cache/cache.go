package cache

import (
	"context"
	"fmt"
	"time"
)

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
	RPop(key string, count *int) any
	LPop(key string, count *int) any
	Type(key string) string
	BLPop(key string, timeout float64) chan any
}
type cache struct {
	data               map[any]any
	listData           map[any][]any
	blockedClients     []chan any
	listDataInsertChan chan struct{}
}

func New() Cache {
	c := &cache{
		data:               make(map[any]any),
		listData:           make(map[any][]any),
		blockedClients:     []chan any{},
		listDataInsertChan: make(chan struct{}, 3),
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
	if len(c.blockedClients) > 0 { // spread event
		c.listDataInsertChan <- struct{}{}
	}

	return len(v) + len(data)
}

func (c *cache) LPush(key string, data []any) int {
	v, _ := c.listData[key]
	c.listData[key] = append(data, v...)
	if len(c.blockedClients) > 0 { // spread event
		c.listDataInsertChan <- struct{}{}
	}

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

func (c *cache) RPop(key string, count *int) any {
	v, _ := c.listData[key]
	if len(v) == 0 {
		return nil
	}

	if count == nil { // default
		r := v[len(v)-1]
		c.listData[key] = v[:len(v)-1]
		return fmt.Sprintf("%v", r)
	}

	nCount := *count
	if nCount > len(v) {
		return v
	}

	endIdx := len(v) - 1 - nCount
	r := v[endIdx:]
	if nCount == len(v) {
		c.listData[key] = nil
		return r
	}

	c.listData[key] = v[:endIdx]
	return r
}

func (c *cache) LPop(key string, count *int) any {
	v, _ := c.listData[key]
	if len(v) == 0 {
		return nil
	}

	if count == nil { // default
		r := v[0]
		c.listData[key] = v[1:]
		return fmt.Sprintf("%v", r)
	}

	nCount := *count
	if nCount > len(v) {
		return v
	}

	r := v[0:nCount]
	if nCount == len(v) {
		c.listData[key] = nil
		return r
	}

	c.listData[key] = v[nCount:]
	return r
}

func (c *cache) Type(key string) string {
	v, ok := c.data[key]
	if !ok {
		v, ok = c.listData[key]
		if !ok {
			return "none"
		}
	}

	switch v.(type) {
	case string:
		return "string"
	case int:
		return "int"
	case bool:
		return "bool"
	case error:
		return "error"
	case []any:
		return "array"
	case map[any]any:
		return "map"
	default:
		return "none"
	}
}
func (c *cache) BLPop(key string, timeout float64) chan any {
	commChan := make(chan any)
	c.blockedClients = append(c.blockedClients, commChan)
	go func() {
		v, ok := c.listData[key]
		if ok || len(v) > 0 {
			c.listData[key] = v[1:]
			commChan <- []any{key, v[0]}
			return
		}

		ctx := context.Background()
		if timeout > 0 {
			ctx, _ = context.WithTimeout(ctx, time.Duration(timeout)*time.Second+2*time.Millisecond)
		}

		for {
			select {
			case <-ctx.Done():
				commChan <- ""
				return
			case <-c.listDataInsertChan:
				firstChan := c.blockedClients[0]
				v, ok := c.listData[key]
				if ok || len(v) > 0 {
					c.listData[key] = v[1:]
					firstChan <- []any{key, v[0]}
					return
				}
			}
		}
	}()
	return commChan
}
