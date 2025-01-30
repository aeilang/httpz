package httpz

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
)

// HelperResponseWriter is a built-in type used to record the StatusCode
// and quickly send responses.
type HelperResponseWriter struct {
	http.ResponseWriter
}

// NewHelperRW creates a new instance of HelperResponseWriter.
func NewHelperRW(w http.ResponseWriter) *HelperResponseWriter {
	return &HelperResponseWriter{
		ResponseWriter: w,
	}
}

// Flush implements the http.Flusher interface.
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

// Hijack implements the http.Hijacker interface.
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

// Push implements the http.Pusher interface.
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

// Unwrap returns the wrapped ResponseWriter.
func (rw *HelperResponseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

// Map is a type alias for a map with string keys and any type values.
type Map map[string]any

// JSON sends a JSON response with the specified status code and data.
func (rw *HelperResponseWriter) JSON(statusCode int, data any) error {
	rw.Header().Set(HeaderContentType, MIMEApplicationJSON)
	rw.WriteHeader(statusCode)
	return json.NewEncoder(rw).Encode(data)
}

// String sends a plain text response with the specified status code and string.
func (rw *HelperResponseWriter) String(statusCode int, s string) error {
	rw.Header().Set(HeaderContentType, MIMETextPlainCharsetUTF8)
	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(s))
	return err
}

// HTML sends an HTML response with the specified status code and HTML content.
func (rw *HelperResponseWriter) HTML(statusCode int, html string) error {
	rw.Header().Set(HeaderContentType, MIMETextHTMLCharsetUTF8)
	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(html))
	return err
}

// XML sends an XML response with the specified status code and data.
func (rw *HelperResponseWriter) XML(statusCode int, data any, indent string) error {
	rw.Header().Set(HeaderContentType, MIMEApplicationXMLCharsetUTF8)
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

// JSON is a convenience function for sending a JSON response.
// the same as
//
//	hw := NewHelperRW(w)
//	return hw.JSON(statusCode, data)
func JSON(w http.ResponseWriter, statusCode int, data any) error {
	hw := NewHelperRW(w)
	return hw.JSON(statusCode, data)
}

// String is a convenience function for sending a plain text response.
// the same as
//
//	hw := NewHelperRW(w)
//	return hw.String(statusCode, s)
func String(w http.ResponseWriter, statusCode int, s string) error {
	hw := NewHelperRW(w)
	return hw.String(statusCode, s)
}

// HTML is a convenience function for sending an HTML response.
// the same as
//
//	hw := NewHelperRW(w)
//	return hw.HTML(statusCode, html)
func HTML(w http.ResponseWriter, statusCode int, html string) error {
	hw := NewHelperRW(w)
	return hw.HTML(statusCode, html)
}

// XML is a convenience function for sending an XML response.
// the same as
//
//	hw := NewHelperRW(w)
//	return hw.XML(statusCode, data, indent)
func XML(w http.ResponseWriter, statusCode int, data any, indent string) error {
	hw := NewHelperRW(w)
	return hw.XML(statusCode, data, indent)
}
