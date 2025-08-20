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

func handleLPush(args []string) (string, error) {
	if len(args) < 2 {
		return protocol.ErrorString("ERR wrong number of arguments for 'lpush' command"), nil
	}

	otherArgs := args[1:]
	anyArgs := make([]any, len(otherArgs))
	for i, a := range otherArgs {
		anyArgs[len(otherArgs)-i-1] = a
	}

	r := cache.LPush(args[0], anyArgs)
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
	if len(r) == 0 {
		return protocol.Array([]any{}), nil
	}

	return protocol.Array(r), nil
}

func handleLLen(args []string) (string, error) {
	if len(args) < 1 {
		return protocol.ErrorString("ERR wrong number of arguments for 'llen' command"), nil
	}
	r := cache.LLen(args[0])
	return protocol.Integer(r), nil
}

func handleRPop(args []string) (string, error) {
	if len(args) < 1 {
		return protocol.ErrorString("ERR wrong number of arguments for 'rpop' command"), nil
	}

	idx, err := extractPopArgs(args)
	if err != nil {
		return protocol.ErrorString("ERR invalid index argument for 'rpop' command"), nil
	}

	r := cache.RPop(args[0], idx)
	switch r.(type) {
	case string:
		return protocol.BulkString(r.(string)), nil
	default:
		d, _ := r.([]any)
		return protocol.Array(d), nil
	}
}

func handleLPop(args []string) (string, error) {
	if len(args) < 1 {
		return protocol.ErrorString("ERR wrong number of arguments for 'lpop' command"), nil
	}

	idx, err := extractPopArgs(args)
	if err != nil {
		return protocol.ErrorString("ERR invalid index argument for 'lpop' command"), nil
	}

	r := cache.LPop(args[0], idx)
	switch r.(type) {
	case string:
		return protocol.BulkString(r.(string)), nil
	default:
		d, _ := r.([]any)
		return protocol.Array(d), nil
	}
}

func extractPopArgs(args []string) (*int, error) {
	var (
		idx *int
		a   int
		err error
	)
	if len(args) > 1 {
		otherArgs := args[1]
		a, err = strconv.Atoi(strings.TrimSpace(otherArgs))
		idx = &a
	}
	return idx, err
}

func handleBLPop(args []string) (string, error) {
	if len(args) < 2 {
		return protocol.ErrorString("ERR wrong number of arguments for 'blpop' command"), nil
	}

	timeout, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		return protocol.ErrorString("ERR invalid timeout argument for 'blpop' command"), nil
	}

	blChan := cache.BLPop(args[0], timeout)
	r := <-blChan
	switch r.(type) {
	case string:
		return protocol.BulkString(r.(string)), nil
	default:
		d, _ := r.([]any)
		return protocol.Array(d), nil
	}
}
