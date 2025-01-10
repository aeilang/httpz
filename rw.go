package httpz

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// 内置类型，仅用来记录StatusCode。
type ResponseWriter struct {
	http.ResponseWriter
	isCommited bool
	statusCode int
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.isCommited = true
	rw.statusCode = statusCode

	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

type Map map[string]any

// 下面是一些helper方法，用于响应数据。
// 响应JSON数据
func (rw *ResponseWriter) JSON(statusCode int, data any) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	return json.NewEncoder(rw).Encode(data)
}

// 响应String
func (rw *ResponseWriter) String(statusCode int, s string) error {
	rw.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(s))
	return err
}

// 响应HTML
func (rw *ResponseWriter) HTML(statusCode int, html string) error {
	rw.Header().Set("Content-Type", "text/html; charset=UTF-8")
	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(html))
	return err
}

// 响应XML
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
