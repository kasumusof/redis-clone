package executor

import (
	"strings"

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
		return handleEcho(args)
	case "ping":
		return handlePing(args)
	case "set":
		return handleSet(args)
	case "get":
		return handleGet(args)
	case "rpush":
		return handleRPush(args)
	case "lpush":
		return handleLPush(args)
	case "lrange":
		return handleLRange(args)
	case "llen":
		return handleLLen(args)
	case "rpop":
		return handleRPop(args)
	case "lpop":
		return handleLPop(args)
	case "blpop":
		return handleBLPop(args)
	case "type":
		return handleType(args)
	case "xadd":
		return handleXAdd(args)
	case "xrange":
		return handleXRange(args)
	case "xread":
		return handleXRead(args)
	default:
		return errorString, nil

	}
}
