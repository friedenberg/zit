package remote_http

import (
	"bytes"
	"io"
	"net/http"
)

type BufferedResponseWriter struct {
	Dirty    bool
	Response http.Response
	Buffer   bytes.Buffer
}

func (responseWriter *BufferedResponseWriter) GetResponseWriter() http.ResponseWriter {
	return responseWriter
}

func (responseWriter *BufferedResponseWriter) Reset() {
	responseWriter.Dirty = false
	responseWriter.Response = http.Response{
		Header: make(http.Header),
	}
	responseWriter.Buffer.Reset()
}

func (responseWriter *BufferedResponseWriter) Header() http.Header {
	responseWriter.Dirty = true
	return responseWriter.Response.Header
}

func (responseWriter *BufferedResponseWriter) WriteHeader(statusCode int) {
	responseWriter.Dirty = true
	responseWriter.Response.StatusCode = statusCode
}

func (responseWriter *BufferedResponseWriter) Write(p []byte) (int, error) {
	responseWriter.Dirty = true
	return responseWriter.Buffer.Write(p)
}

func (responseWriter *BufferedResponseWriter) WriteResponse(
	writer io.Writer,
) error {
	responseWriter.Response.Body = io.NopCloser(&responseWriter.Buffer)
	responseWriter.Response.ContentLength = int64(responseWriter.Buffer.Len())
	return responseWriter.Response.Write(writer)
}
