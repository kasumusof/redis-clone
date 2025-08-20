package executor

import (
	"github.com/codecrafters-io/redis-starter-go/app/cache"
	"github.com/codecrafters-io/redis-starter-go/app/protocol"
)

func handleGet(args []string) (string, error) {
	if len(args) < 1 {
		return protocol.ErrorString("ERR wrong number of arguments for 'get' command"), nil
	}
	val, ok := cache.Get(args[0])
	if !ok {
		return protocol.Nulls(), nil
	}
	return protocol.BulkString(val.(string)), nil
}
