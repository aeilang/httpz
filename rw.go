package httpz

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
)

// A built-in type, used only to record the StatusCode
// and quickly send responses.
type ResponseWriter struct {
	http.ResponseWriter
	http.Pusher
	isCommited bool
	statusCode int
}

// rewrite the WriteHeader method to record statusCode
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.isCommited = true
	rw.statusCode = statusCode

	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

// implement http.Flusher
func (rw *ResponseWriter) Flush() {
	w := rw.ResponseWriter

	for {
		switch t := w.(type) {
		case http.Flusher:
			t.Flush()
			return
		case rwUnwrapper:
			w = t.Unwrap()
		default:
			return
		}
	}
}

// implement http.Hijacker
func (rw *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	w := rw.ResponseWriter
	for {
		switch t := w.(type) {
		case http.Hijacker:
			return t.Hijack()
		case rwUnwrapper:
			w = t.Unwrap()
		default:
			return nil, nil, fmt.Errorf("the ResponseWriter didn't implement http.Hijacker")
		}
	}
}

// implement http.Pusher
func (rw *ResponseWriter) Push(target string, opts *http.PushOptions) error {
	w := rw.ResponseWriter

	for {
		switch t := w.(type) {
		case http.Pusher:
			return t.Push(target, opts)
		case rwUnwrapper:
			w = t.Unwrap()
		default:
			return fmt.Errorf("the ResponseWriter didn't implement http.Pusher")
		}
	}
}

type rwUnwrapper interface {
	Unwrap() http.ResponseWriter
}

// get the wrapped ResponseWriter
func (rw *ResponseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

type Map map[string]any

// send json
func (rw *ResponseWriter) JSON(statusCode int, data any) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	return json.NewEncoder(rw).Encode(data)
}

// send string
func (rw *ResponseWriter) String(statusCode int, s string) error {
	rw.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(s))
	return err
}

// send html
func (rw *ResponseWriter) HTML(statusCode int, html string) error {
	rw.Header().Set("Content-Type", "text/html; charset=UTF-8")
	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(html))
	return err
}

// send xml
func (rw *ResponseWriter) XML(statusCode int, data any, indent string) error {
	rw.Header().Set("Content-Type", "application/xml; charset=UTF-8")
	rw.WriteHeader(statusCode)
	enc := xml.NewEncoder(rw)
	if indent != "" {
		enc.Indent("", indent)
	}
	if _, err := rw.Write([]byte(xml.Header)); err != nil {
		return err
	}

	return enc.Encode(data)
}
