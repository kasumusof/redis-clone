package cache

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
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
	XAdd(key string, id string, elems []any) (string, bool)
	XRange(key string, start string, end string) []any
	XRead(key string, id string) []any
}
type cache struct {
	data               map[any]any
	listData           map[any][]any
	blockedClients     []chan any
	listDataInsertChan chan struct{}
	streamData         map[any][][2]any
}

func New() Cache {
	c := &cache{
		data:               make(map[any]any),
		listData:           make(map[any][]any),
		blockedClients:     []chan any{},
		listDataInsertChan: make(chan struct{}, 3),
		streamData:         make(map[any][][2]any),
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
			v, ok = c.streamData[key]
			if !ok {
				return "none"
			}
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
	case [][2]any:
		return "stream"
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
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout*float64(time.Second)))
			defer cancel()
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

func (c *cache) XAdd(key string, id string, elems []any) (string, bool) {
	m, ok := c.streamData[key]
	if !ok {
		m = make([][2]any, 0)
	}

	var valid bool
	id, valid = validateXAddID(id, m)
	if !valid {
		return id, false
	}

	d := [2]any{id, elems}
	m = append(m, d)

	c.streamData[key] = m
	return id, true
}

func validateXAddID(id string, data [][2]any) (string, bool) {

	if id == "*" {
		id = fmt.Sprintf("%d-*", time.Now().UnixMilli())
	}

	if len(data) == 0 {
		if strings.HasSuffix(id, "-*") {
			firstSplit := strings.Split(id, "-")[0]
			id = "0-1"
			if firstSplit > "0" {
				id = firstSplit + "-0"
			}
		}

		return id, true
	}

	last := data[len(data)-1]

	idSplit := strings.Split(id, "-")
	if len(idSplit) != 2 {
		return id, false
	}

	newTime, newIncr := idSplit[0], idSplit[1]
	lastSplit := strings.Split(last[0].(string), "-")
	if len(lastSplit) != 2 {
		return id, false
	}

	lastTime, lastIncr := lastSplit[0], lastSplit[1]

	switch {
	case newTime < lastTime:
		return id, false
	case newIncr == "*":
		toAdd := "0"
		if newTime == lastTime {
			lastIncr, _ := strconv.Atoi(lastIncr)
			toAdd = strconv.Itoa(lastIncr + 1)
		}

		id = newTime + "-" + toAdd
		return id, true
	case newTime == lastTime && newIncr <= lastIncr:
		return "", false
	default:
		return id, true
	}
}

func (c *cache) XRange(key string, start string, end string) []any {
	x, ok := c.streamData[key]
	if !ok {
		return nil
	}

	start, end = validateXRangeFilters(start, end)

	var res []any
	for _, v := range x {
		id := v[0].(string)
		if idISGreaterOrEqual(id, start) && idIsLessOrEqual(id, end) {
			res = append(res, v)
		}
	}
	return res
}

func validateXRangeFilters(start string, end string) (string, string) {
	if start == "-" {
		start = "0-1"
	} else if len(strings.Split(start, "-")) == 1 {
		start = start + "-0"
	}

	if end == "+" {
		end = fmt.Sprintf("%d-*", math.MaxInt64)
	} else if len(strings.Split(end, "-")) == 1 {
		end = end + "-0"
	}

	return start, end
}

func idISGreaterOrEqual(main string, target string) bool {
	timeMain, incrMain := strings.Split(main, "-")[0], strings.Split(main, "-")[1]
	timeTarget, incrTarget := strings.Split(target, "-")[0], strings.Split(target, "-")[1]
	if timeMain > timeTarget {
		return true
	} else if timeMain == timeTarget {
		incrMain, _ := strconv.Atoi(incrMain)
		incrTarget, _ := strconv.Atoi(incrTarget)
		return incrMain >= incrTarget
	}
	return false
}

func idIsLessOrEqual(main string, target string) bool {
	timeMain, incrMain := strings.Split(main, "-")[0], strings.Split(main, "-")[1]
	timeTarget, incrTarget := strings.Split(target, "-")[0], strings.Split(target, "-")[1]
	if timeMain < timeTarget {
		return true
	} else if timeMain == timeTarget {
		incrMain, _ := strconv.Atoi(incrMain)
		incrTarget, _ := strconv.Atoi(incrTarget)
		return incrMain <= incrTarget
	}
	return false
}

func idIsGreater(main string, target string) bool {
	timeMain, incrMain := strings.Split(main, "-")[0], strings.Split(main, "-")[1]
	timeTarget, incrTarget := strings.Split(target, "-")[0], strings.Split(target, "-")[1]
	if timeMain > timeTarget {
		return true
	} else if timeMain == timeTarget {
		incrMain, _ := strconv.Atoi(incrMain)
		incrTarget, _ := strconv.Atoi(incrTarget)
		return incrMain > incrTarget
	}
	return false
}

func (c *cache) XRead(key string, idTarget string) []any {
	var res []any
	v, ok := c.streamData[key]
	if !ok {
		return res
	}

	for _, v := range v {
		id := v[0].(string)
		if idIsGreater(id, idTarget) {
			res = append(res, v)
			continue
		}
	}

	res = append([]any{key}, res)
	return []any{res}
}
