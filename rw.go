package httpz

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// A built-in type, used only to record the StatusCode
// and quickly send responses.
type ResponseWriter struct {
	http.ResponseWriter
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

// Is this necessary?
func (rw *ResponseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

type Map map[string]any

// response json
func (rw *ResponseWriter) JSON(statusCode int, data any) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	return json.NewEncoder(rw).Encode(data)
}

// string
func (rw *ResponseWriter) String(statusCode int, s string) error {
	rw.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(s))
	return err
}

// html
func (rw *ResponseWriter) HTML(statusCode int, html string) error {
	rw.Header().Set("Content-Type", "text/html; charset=UTF-8")
	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(html))
	return err
}

// xml
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
