package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	connChan := readConnection(l)

	parseCommand(connChan)

}

func readConnection(l net.Listener) chan net.Conn {
	connChan := make(chan net.Conn)

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				os.Exit(1)
			}

			connChan <- conn
		}
	}()
	return connChan
}

func parseCommand(connChan chan net.Conn) {
	for {
		select {
		case conn := <-connChan:
			func(conn net.Conn) {
				defer conn.Close()

				data := make([]byte, 0)
				if _, err := conn.Read(data); err != nil {
					fmt.Println("Error reading request: ", err.Error())
					return
				}

				command := string(data)
				resp := "+PONG\r\n"
				for range strings.Split(command, "\n") {
					resp += "+PONG\r\n"
				}
				_, _ = conn.Write([]byte(resp))
			}(conn)
		}
	}
}
