package executor

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/protocol"
)

func handlePing(args []string) (string, error) {
	if len(args) > 0 {
		return protocol.ErrorString(strings.Join(args, " ")), nil
	}
	return protocol.SimpleString("PONG"), nil
}
