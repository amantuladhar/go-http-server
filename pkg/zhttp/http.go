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
	method string
	rg     *regexp.Regexp
	fn     handleFunc
}

type httpServerConfig struct {
	paths map[string]pathDetail
}

func (c *httpServerConfig) HandleFunc(fullPattern string, handler handleFunc) error {
	// Split method and path
	firstSplit := strings.SplitN(fullPattern, " ", 2)
	if len(firstSplit) != 2 {
		return fmt.Errorf("pattern on handle func must be <HTTP_VERB><SPACE><ROUTE_PATH>. e.g `GET /api/v1/yellos`")
	}

	split := strings.Split(firstSplit[1], "/")
	splitPattern := make([]string, 0, len(firstSplit[1]))
	for _, s := range split {
		if strings.Contains(s, " ") {
			return fmt.Errorf("path placeholder cannot have whitespaces")
		}
		if len(s) > 3 && s[0] == '{' && s[len(s)-1] == '}' {
			// This is better but cc is only considers 2 level paths
			// "(?<%s>[\\w]*[^\\/])"
			splitPattern = append(splitPattern, fmt.Sprintf("(?<%s>.*)", s[1:len(s)-1]))
		} else {
			splitPattern = append(splitPattern, s)
		}
	}

	finalPattern := fmt.Sprintf("^%s$", strings.Join(splitPattern, "/"))
	rg, err := regexp.Compile(finalPattern)
	if err != nil {
		util.LogErr("unable to compile regex", "path", fullPattern, "pattern", finalPattern, "err", err)
		return nil
	}
	c.paths[fullPattern] = pathDetail{
		method: firstSplit[0],
		rg:     rg,
		fn:     handler,
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
		go func() {
			req, err := parseRequest(&conn)
			if err != nil {
				util.LogErr("unable to parse request object from client connection", "err", err)
				return
			}
			err = processRequest(&conn, config, req)
			if err != nil {
				util.LogErr("unable to process client request", "err", err)
				return
			}

			conn.Close()
		}()
	}
}

func processRequest(conn *net.Conn, config *httpServerConfig, req *Request) error {
	response := NewResponse().StatusCode(404)

	for _, detail := range config.paths {
		if req.Method != detail.method || !detail.rg.MatchString(req.Path) {
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

	headerList := make([]string, 0, len(response.Headers))
	for k, v := range response.Headers {
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
	case 201:
		return "Created"
	case 404:
		return "Not Found"
	default:
		return "Internal Server Error"
	}
}

func parseRequest(conn *net.Conn) (*Request, error) {
	// TODO: Update to use buffered reader way
	// Also what if 1024 bytes is not enought
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

	headers := make(map[string]string)
	// incrementing this count so we know where body part is
	lineCount := 1
	var prevLine string
	for _, line := range lines[1:] {
		lineCount += 1
		if prevLine == "" && line == "" {
			break
		}
		splitH := strings.Split(strings.TrimSpace(line), ": ")
		headers[splitH[0]] = splitH[1]
	}

	body := []byte(lines[lineCount])

	return &Request{
		Method:      method,
		Path:        path,
		HttpVersion: httpVersion,
		PathParam:   make(map[string]string),
		Headers:     headers,
		Body:        body,
	}, nil
}

func parseStatusLine(line string) (method string, path string, version string, err error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid status line")
	}
	return parts[0], parts[1], parts[2], nil
}
