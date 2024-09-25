package main

import (
	"fmt"
	"log/slog"
	"net"
	"regexp"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/pkg/util"
)

type HandleFunc func(*Request) *Response

type HttpServerConfig struct {
	paths map[string]HandleFunc
}

func (c *HttpServerConfig) HandleFunc(pattern string, handler HandleFunc) *HttpServerConfig {
	c.paths[pattern] = handler
	return c
}

func NewHttpServerConfig() *HttpServerConfig {
	return &HttpServerConfig{
		paths: make(map[string]HandleFunc),
	}
}

func ListenAndServe(addr string, config *HttpServerConfig) error {
	l, err := net.Listen("tcp", addr)
	util.ExitOnErr(err, "failed to bind", "addr", addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			util.LogErr("unable to accept client connection", "err", err)
			continue
		}

		req, err := parseRequest(&conn)
		if err != nil {
			util.LogErr("unable to parse request object from client connection", "err", err)
			continue
		}

		err = processRequest(&conn, config, req)
		if err != nil {
			util.LogErr("unable to process client request", "err", err)
			continue
		}

		conn.Close()
	}
}

type Response struct {
	headers    map[string]string
	body       []byte
	statusCode int
}

func NewResponse() *Response {
	return &Response{
		statusCode: 200,
	}
}

func (c *Response) Json(body []byte) *Response {
	c.headers["Content-Type"] = "application/json"
	c.body = body
	return c
}
func (r *Response) StatusCode(code int) *Response {
	r.statusCode = code
	return r
}

func processRequest(conn *net.Conn, config *HttpServerConfig, req *Request) error {
	response := NewResponse().StatusCode(404)

	// Need to update to use regex
	for path, fn := range config.paths {
		split := strings.Split(path, "/")
		splitPattern := make([]string, 0, len(split))
		for _, s := range split {
			if strings.Contains(s, " ") {
				return fmt.Errorf("path placeholder cannot have whitespaces")
			}
			if len(s) > 3 && s[0] == '{' && s[len(s)-1] == '}' {
				// This is better but cc is only considering 1 level paths
				// "(?<%s>[\\w]*[^\\/])"
				splitPattern = append(splitPattern, fmt.Sprintf("(?<%s>.*)", s[1:len(s)-1]))
			} else {
				splitPattern = append(splitPattern, s)
			}
		}
		finalPattern := fmt.Sprintf("^%s$", strings.Join(splitPattern, "/"))
		rg, err := regexp.Compile(finalPattern)
		util.LogOnErr(err, "unable to compile regex", "path", path, "pattern", "pattern")

		if !rg.MatchString(req.Path) {
			continue
		}
		// fmt.Println(re.MatchString(xx))
		// matches := re.FindStringSubmatch(xx)
		// fmt.Println(fmt.Sprintf("0 -- %s", matches))
		// fmt.Printf("1. %q\n", re.SubexpNames())

		grpNames := rg.SubexpNames()
		values := rg.FindStringSubmatch(path)
		for i := 1; i < len(grpNames); i++ {
			req.PathParam[grpNames[i]] = values[i]
		}
		response = fn(req)
		break
	}

	// build response string
	var statusLine = fmt.Sprintf("HTTP/1.1 %d %s", response.statusCode, statusCodeString(response.statusCode))

	var headers = ""
	var body = ""

	var finalResponse = fmt.Sprintf("%s\r\n%s\r\n%s", statusLine, headers, body)

	_, err := (*conn).Write([]byte(finalResponse))
	return err
}

func statusCodeString(code int) string {
	switch code {
	case 200:
		return "OK"
	case 404:
		return "Not Found"
	default:
		return "Internal Server Error"
	}
}

func parseRequest(conn *net.Conn) (*Request, error) {
	// TODO: Need to find a buffered reader way to do this
	buf := make([]byte, 1024)
	n, err := (*conn).Read(buf)
	if err != nil {
		return nil, fmt.Errorf("unable to read from client connection")
	}
	lines := strings.Split(string(buf[:n]), "\r\n")
	method, path, httpVersion, err := parseStatusLine(lines[0])
	if err != nil {
		return nil, fmt.Errorf("unable to parse status line")
	}
	return &Request{
		Method:      method,
		Path:        path,
		HttpVersion: httpVersion,
	}, nil
}

type Request struct {
	Method      string
	Path        string
	HttpVersion string
	PathParam   map[string]string
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	config := NewHttpServerConfig()
	config.HandleFunc("/", func(r *Request) *Response {
		return NewResponse().StatusCode(200)
	})
	config.HandleFunc("/echo/{slug}", func(r *Request) *Response {
		return NewResponse().StatusCode(200)
	})
	err := ListenAndServe("0.0.0.0:4221", config)
	util.ExitOnErr(err, "unable to start server")
}

func parseStatusLine(line string) (method string, path string, version string, err error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid status line")
	}
	return parts[0], parts[1], parts[2], nil
}
