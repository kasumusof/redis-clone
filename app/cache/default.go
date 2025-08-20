package cache

import "time"

var (
	defaultCache Cache
)

func init() {
	defaultCache = New()
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

func RPush(key string, value []any) int {
	return defaultCache.RPush(key, value)
}

func LRange(key string, start, end int) []any {
	return defaultCache.LRange(key, start, end)
}

func LPush(key string, args []any) int {
	return defaultCache.LPush(key, args)
}

func LLen(key string) int {
	return defaultCache.LLen(key)
}

func RPop(key string, count *int) any {
	return defaultCache.RPop(key, count)
}

func LPop(key string, count *int) any {
	return defaultCache.LPop(key, count)
}

func Type(key string) string {
	return defaultCache.Type(key)
}

func BLPop(s string, timeout float64) chan any {
	return defaultCache.BLPop(s, timeout)
}

func XAdd(key string, id string, elems []any) (string, bool) {
	return defaultCache.XAdd(key, id, elems)
}

func XRange(key, start, end string) []any {
	return defaultCache.XRange(key, start, end)
}

func XRead(key string, id string) []any {
	return defaultCache.XRead(key, id)
}
