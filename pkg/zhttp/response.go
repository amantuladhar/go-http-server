package zhttp

import "fmt"

type Response struct {
	headers    map[string]string
	body       []byte
	statusCode int
}

func NewResponse() *Response {
	return &Response{
		headers:    make(map[string]string),
		statusCode: 200,
	}
}

func (c *Response) setBody(body []byte, contentType string) {
	c.headers["Content-Type"] = contentType
	c.headers["Content-Length"] = fmt.Sprintf("%d", len(body))
	c.body = body
}

func (c *Response) Json(body []byte) *Response {
	c.setBody(body, "application/json")
	return c
}

func (c *Response) Text(body []byte) *Response {
	c.setBody(body, "text/plain")
	return c
}

func (r *Response) StatusCode(code int) *Response {
	r.statusCode = code
	return r
}
