package executor

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/cache"
	"github.com/codecrafters-io/redis-starter-go/app/protocol"
)

func handleXAdd(args []string) (string, error) {
	if len(args) < 2 {
		return protocol.ErrorString("ERR wrong number of arguments for 'xadd' command"), nil
	}

	key := args[0]
	id := args[1]

	if id == "0-0" {
		return protocol.ErrorString("ERR The ID specified in XADD must be greater than 0-0"), nil
	}

	if len(args) < 3 || len(args)%2 != 0 {
		return protocol.ErrorString("ERR wrong number of arguments for 'xadd' command"), nil
	}

	otherArgs := make([]any, len(args[2:]))
	for i := 0; i < len(otherArgs); i++ {
		otherArgs[i] = args[i+2]
	}

	r, ok := cache.XAdd(key, id, otherArgs)
	if !ok {
		return protocol.ErrorString("ERR The ID specified in XADD is equal or smaller than the target stream top item"), nil
	}

	return protocol.BulkString(r), nil
}

func handleXRange(args []string) (string, error) {
	if len(args) < 3 {
		return protocol.ErrorString("ERR wrong number of arguments for 'xrange' command"), nil
	}

	start := strings.TrimSpace(args[1])
	end := strings.TrimSpace(args[2])

	r := cache.XRange(args[0], start, end)
	return protocol.Array(r), nil
}
