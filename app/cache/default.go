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

func RPop(key string) string {
	return defaultCache.RPop(key)
}

func LPop(key string) string {
	return defaultCache.LPop(key)
}
