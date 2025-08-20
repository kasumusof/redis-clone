package executor

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/cache"
	"github.com/codecrafters-io/redis-starter-go/app/protocol"
)

func handleRPush(args []string) (string, error) {
	if len(args) < 2 {
		return protocol.ErrorString("ERR wrong number of arguments for 'rpush' command"), nil
	}

	anyArgs := make([]any, len(args[1:]))
	for i, a := range args[1:] {
		anyArgs[i] = a
	}

	r := cache.RPush(args[0], anyArgs)
	return protocol.Integer(r), nil
}

func handleLRange(args []string) (string, error) {
	if len(args) < 3 {
		return protocol.ErrorString("ERR wrong number of arguments for 'lrange' command"), nil
	}

	start, err := strconv.Atoi(strings.TrimSpace(args[1]))
	if err != nil {
		return protocol.ErrorString("ERR invalid start argument for 'lrange' command"), nil
	}
	end, err := strconv.Atoi(strings.TrimSpace(args[2]))
	if err != nil {
		return protocol.ErrorString("ERR invalid end argument for 'lrange' command"), nil
	}

	r := cache.LRange(args[0], start, end)
	return protocol.Array(r), nil
}
