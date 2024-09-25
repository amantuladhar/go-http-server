package zhttp

import "fmt"

type Response struct {
	Headers    map[string]string
	body       []byte
	statusCode int
}

func NewResponse() *Response {
	return &Response{
		Headers:    make(map[string]string),
		statusCode: 200,
	}
}

func (c *Response) setBody(body []byte, contentType string) {
	c.Headers["Content-Type"] = contentType
	c.Headers["Content-Length"] = fmt.Sprintf("%d", len(body))
	c.body = body
}

func (c *Response) Json(body []byte) *Response {
	c.setBody(body, "application/json")
	return c
}

func (c *Response) File(body []byte) *Response {
	c.setBody(body, "application/octet-stream")
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
