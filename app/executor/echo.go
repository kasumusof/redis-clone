package executor

import "github.com/codecrafters-io/redis-starter-go/app/protocol"

func handleEcho(args []string) (string, error) {
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
}
