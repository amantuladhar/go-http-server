package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/pkg/util"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	fmt.Println("Logs from your program will appear here!")
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	util.ExitOnErr(err, "Failed to bind to port 4221")

	conn, err := l.Accept()
	util.ExitOnErr(err, "Error accepting connection")
	defer conn.Close()

	// parse request
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	util.LogOnErr(err, "unable to read from connection")

	lines := strings.Split(string(buf[:n]), "\r\n")
	_, path, _, err := parseStatusLine(lines[0])
	util.LogOnErr(err, "unable to parse status line")

	response := ""
	if path == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	//response
	conn.Write([]byte(response))
}

func parseStatusLine(line string) (method string, path string, version string, err error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid status line")
	}
	return parts[0], parts[1], parts[2], nil
}
