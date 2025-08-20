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
		return handleEcho(args)
	case "ping":
		return handlePing(args)
	case "set":
		return handleSet(args)
	case "get":
		return handleGet(args)
	default:
		return errorString, nil

	}
}
