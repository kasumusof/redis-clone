package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/executor"
	"github.com/codecrafters-io/redis-starter-go/app/protocol"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var (
	_ = net.Listen
	_ = os.Exit
)

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
			buff := bufio.NewReader(conn)
			go func(conn net.Conn) {
				defer closeConnection(conn)
				for {
					res, err := protocol.ParseRequest(buff)
					if err != nil {
						if errors.Is(err, io.EOF) {
							continue
						}

						log.Println(err)
					}

					resp, err := executor.Execute(res)
					if err != nil {
						log.Println(err)
					}

					if _, err = conn.Write([]byte(resp)); err != nil {
						log.Println(err)
					}
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
