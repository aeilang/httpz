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
type HelperResponseWriter struct {
	http.ResponseWriter
}

func NewHelperRW(w http.ResponseWriter) *HelperResponseWriter {
	return &HelperResponseWriter{
		ResponseWriter: w,
	}
}

// implement http.Flusher
func (rw *HelperResponseWriter) Flush() {
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
func (rw *HelperResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
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
func (rw *HelperResponseWriter) Push(target string, opts *http.PushOptions) error {
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
func (rw *HelperResponseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

type Map map[string]any

// send json
func (rw *HelperResponseWriter) JSON(statusCode int, data any) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	return json.NewEncoder(rw).Encode(data)
}

// send string
func (rw *HelperResponseWriter) String(statusCode int, s string) error {
	rw.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(s))
	return err
}

// send html
func (rw *HelperResponseWriter) HTML(statusCode int, html string) error {
	rw.Header().Set("Content-Type", "text/html; charset=UTF-8")
	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(html))
	return err
}

// send xml
func (rw *HelperResponseWriter) XML(statusCode int, data any, indent string) error {
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
