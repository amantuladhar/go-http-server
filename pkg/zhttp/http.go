package zhttp

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/pkg/util"
)

func NewHttpServerConfig() *httpServerConfig {
	return &httpServerConfig{
		paths: make(map[string]pathDetail),
	}
}

type handleFunc func(*Request) *Response
type pathDetail struct {
	rg *regexp.Regexp
	fn handleFunc
}

type httpServerConfig struct {
	paths map[string]pathDetail
}

func (c *httpServerConfig) HandleFunc(pattern string, handler handleFunc) error {
	split := strings.Split(pattern, "/")
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
	util.LogOnErr(err, "unable to compile regex", "path", pattern, "pattern", finalPattern)
	c.paths[pattern] = pathDetail{
		rg: rg,
		fn: handler,
	}
	return nil
}

func ListenAndServe(addr string, config *httpServerConfig) error {
	l, err := net.Listen("tcp", addr)
	util.ExitOnErr(err, "failed to bind", "addr", addr)
	util.LogInfo(fmt.Sprintf("Server started: http://%s", addr))
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

func processRequest(conn *net.Conn, config *httpServerConfig, req *Request) error {
	response := NewResponse().StatusCode(404)

	for _, detail := range config.paths {
		if !detail.rg.MatchString(req.Path) {
			continue
		}
		grpNames := detail.rg.SubexpNames()
		values := detail.rg.FindStringSubmatch(req.Path)
		for i := 1; i < len(grpNames); i++ {
			req.PathParam[grpNames[i]] = values[i]
		}
		response = detail.fn(req)
		break
	}

	// build response string
	var statusLine = fmt.Sprintf("HTTP/1.1 %d %s", response.statusCode, statusCodeString(response.statusCode))

	headerList := make([]string, 0, len(response.headers))
	for k, v := range response.headers {
		headerList = append(headerList, fmt.Sprintf("%s: %s", k, v))
	}
	var headers = strings.Join(headerList, "\r\n")

	var body = string(response.body)

	var finalResponse = fmt.Sprintf("%s\r\n%s\r\n\r\n%s", statusLine, headers, body)

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
		PathParam:   make(map[string]string),
	}, nil
}

func parseStatusLine(line string) (method string, path string, version string, err error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid status line")
	}
	return parts[0], parts[1], parts[2], nil
}
