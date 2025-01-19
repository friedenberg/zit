package remote_http

import (
	"io"
	"net/http"
	"strings"
)

type MethodPath struct {
	Method string
	Path   string
}

type Request struct {
	MethodPath
	Headers http.Header
	Body    io.ReadCloser
}

type Response struct {
	StatusCode int
	Body       io.ReadCloser
}

func (r *Response) ErrorWithStatus(status int, err error) {
	r.StatusCode = status
	r.Body = io.NopCloser(strings.NewReader(err.Error()))
}

func (r *Response) Error(err error) {
	r.ErrorWithStatus(http.StatusInternalServerError, err)
}
