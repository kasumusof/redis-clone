package executor

import (
	"github.com/codecrafters-io/redis-starter-go/app/cache"
	"github.com/codecrafters-io/redis-starter-go/app/protocol"
)

func handleType(args []string) (string, error) {
	if len(args) < 1 {
		return protocol.ErrorString("ERR wrong number of arguments for 'type' command"), nil
	}

	r := cache.Type(args[0])

	return protocol.SimpleString(r), nil
}
