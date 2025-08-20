package executor

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/cache"
	"github.com/codecrafters-io/redis-starter-go/app/protocol"
)

var (
	errorString = protocol.ErrorString("ERR unknown command")
)

func Execute(resp protocol.RESP) (string, error) {
	var (
		cmd  string
		args []string
	)

	if resp == nil {
		return errorString, nil
	}

	if delim, ok := resp.IsArray(); ok {
		raw := strings.Split(resp.String(), delim)
		if len(raw) > 0 {
			cmd = raw[0]
		}
		if len(raw) > 1 {
			args = raw[1:]
		}
	} else {
		cmd = resp.String()
	}

	switch strings.ToLower(cmd) {
	case "echo":
		if len(args) == 0 {
			return protocol.ErrorString("ERR wrong number of arguments for 'echo' command"), nil
		}

		if len(args) == 1 {
			return protocol.BulkString(args[0]), nil
		}

		argsToUse := make([]any, len(args))
		for a := range args {
			argsToUse[a] = args[a]
		}

		return protocol.Array(argsToUse), nil
	case "ping":
		if len(args) > 0 {
			return protocol.ErrorString(strings.Join(args, " ")), nil
		}
		return protocol.SimpleString("PONG"), nil
	case "set":
		if len(args) < 2 {
			return protocol.ErrorString("ERR wrong number of arguments for 'set' command"), nil
		}
		exArgKey := "ex"
		exExists := false
		var expArg string
		for i := 2; i < len(args); i++ {
			if strings.ToLower(args[i]) == exArgKey {
				exExists = true
				if i+1 < len(args) {
					expArg = args[i+1]
				}
				break
			}
		}

		if expArg == "" && exExists {
			return protocol.ErrorString("ERR wrong number of arguments for 'set' command"), nil
		}

		ex, err := strconv.Atoi(strings.TrimSpace(expArg))
		if err != nil {
			return protocol.ErrorString("ERR wrong number of arguments for 'set' command"), nil
		}

		cache.Set(args[0], args[1], ex)
		return protocol.SimpleString("OK"), nil
	case "get":
		if len(args) < 1 {
			return protocol.ErrorString("ERR wrong number of arguments for 'get' command"), nil
		}
		val, ok := cache.Get(args[0])
		if !ok {
			return protocol.Nulls(), nil
		}
		return protocol.BulkString(val.(string)), nil
	default:
		return errorString, nil

	}
}
