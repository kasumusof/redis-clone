package executor

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/cache"
	"github.com/codecrafters-io/redis-starter-go/app/protocol"
)

func handleGet(args []string) (string, error) {
	if len(args) < 1 {
		return protocol.ErrorString("ERR wrong number of arguments for 'get' command"), nil
	}
	val, ok := cache.Get(args[0])
	if !ok {
		return protocol.BulkString(""), nil
	}
	return protocol.BulkString(val.(string)), nil
}

func handleSet(args []string) (string, error) {
	if len(args) < 2 {
		return protocol.ErrorString("ERR wrong number of arguments for 'set' command"), nil
	}
	exArgKey := "px"
	var expArg string
	for i := 2; i < len(args); i++ {
		if strings.ToLower(args[i]) == exArgKey {
			if i+1 < len(args) {
				expArg = args[i+1]
			}
			break
		}
	}

	var (
		err error
		ex  int
	)

	if expArg != "" {
		ex, err = strconv.Atoi(strings.TrimSpace(expArg))
		if err != nil {
			return protocol.ErrorString("ERR wrong number of arguments for 'set' command"), nil
		}
	}

	cache.Set(args[0], args[1], ex)
	return protocol.SimpleString("OK"), nil
}
