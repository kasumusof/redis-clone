package executor

import (
	"github.com/codecrafters-io/redis-starter-go/app/cache"
	"github.com/codecrafters-io/redis-starter-go/app/protocol"
)

func handleRPush(args []string) (string, error) {
	if len(args) < 2 {
		return protocol.ErrorString("ERR wrong number of arguments for 'rpush' command"), nil
	}
	r := cache.RPush(args[0], args[1])
	return protocol.Integer(r), nil
}
