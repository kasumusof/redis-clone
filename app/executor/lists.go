package executor

import (
	"github.com/codecrafters-io/redis-starter-go/app/cache"
	"github.com/codecrafters-io/redis-starter-go/app/protocol"
)

func handleRPush(args []string) (string, error) {
	if len(args) < 2 {
		return protocol.ErrorString("ERR wrong number of arguments for 'rpush' command"), nil
	}

	anyArgs := make([]any, len(args))
	for i, a := range args[1:] {
		anyArgs[i] = a
	}

	r := cache.RPush(args[0], anyArgs[1:])
	return protocol.Integer(r), nil
}
