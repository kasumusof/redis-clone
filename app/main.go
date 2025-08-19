package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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

	defer func() {
		if err := l.Close(); err != nil {
			log.Fatal(err)
		}
	}()

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
			go func(conn net.Conn) {
				defer closeConnection(conn)
				for {
					buff := bufio.NewReader(conn)
					command, err := buff.ReadString('\n')
					if err != nil {
						if errors.Is(err, io.EOF) {
							continue
						}
						log.Fatal("buffer error", err)
					}

					_ = command

					resp := "+PONG\r\n"
					_, _ = conn.Write([]byte(resp))
				}
			}(conn)
		}
	}
}

func closeConnection(conn net.Conn) {
	if err := conn.Close(); err != nil {
		log.Fatal(err)
	}
}
