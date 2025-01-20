package remote_http

import (
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type MethodPath struct {
	Method string
	Path   string
}

type Request struct {
	request *http.Request
	MethodPath
	Headers http.Header
	Body    io.ReadCloser
}

func (r Request) Vars() map[string]string {
	return mux.Vars(r.request)
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
