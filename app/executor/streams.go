package executor

import (
	"github.com/codecrafters-io/redis-starter-go/app/cache"
	"github.com/codecrafters-io/redis-starter-go/app/protocol"
)

func handleXAdd(args []string) (string, error) {
	if len(args) < 2 {
		return protocol.ErrorString("ERR wrong number of arguments for 'xadd' command"), nil
	}

	key := args[0]
	id := args[1]

	if len(args) < 3 || len(args)%2 != 0 {
		return protocol.ErrorString("ERR wrong number of arguments for 'xadd' command"), nil
	}

	otherArgs := make([][2]any, len(args[2:])/2)
	for i := 0; i < len(otherArgs); i++ {
		otherArgs[i] = [2]any{args[2+i*2], args[2+i*2+1]}
	}

	r := cache.XAdd(key, id, otherArgs)
	return protocol.BulkString(r), nil
}
