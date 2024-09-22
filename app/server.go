package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/pkg/util"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	fmt.Println("Logs from your program will appear here!")
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	util.ExitOnErr(err, "Failed to bind to port 4221")

	conn, err := l.Accept()
	util.ExitOnErr(err, "Error accepting connection")
	defer conn.Close()
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}
