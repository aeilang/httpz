package httpz

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelperResponseWriter_JSON(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewHelperRW(rec)

	data := Map{"key": "value"}
	err := rw.JSON(http.StatusOK, data)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	expectedBody := `{"key":"value"}`
	assert.JSONEq(t, expectedBody, rec.Body.String())
	assert.Equal(t, MIMEApplicationJSON, rec.Header().Get(HeaderContentType))
}

func TestHelperResponseWriter_String(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewHelperRW(rec)

	err := rw.String(http.StatusOK, "hello world")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	expectedBody := "hello world"
	assert.Equal(t, expectedBody, rec.Body.String())
	assert.Equal(t, MIMETextPlainCharsetUTF8, rec.Header().Get(HeaderContentType))
}

func TestHelperResponseWriter_HTML(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewHelperRW(rec)

	err := rw.HTML(http.StatusOK, "<h1>Hello</h1>")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	expectedBody := "<h1>Hello</h1>"
	assert.Equal(t, expectedBody, rec.Body.String())
	assert.Equal(t, MIMETextHTMLCharsetUTF8, rec.Header().Get(HeaderContentType))
}

type User struct {
	Name string `xml:"name"`
	Age  int    `xml:"age"`
}

func TestHelperResponseWriter_XML(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewHelperRW(rec)

	u := User{
		Name: "lihua",
		Age:  18,
	}
	err := rw.XML(http.StatusOK, u, "  ")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	expectedBody := `<?xml version="1.0" encoding="UTF-8"?>` + "\n" + "<User>\n  <name>lihua</name>\n  <age>18</age>\n</User>"
	assert.Equal(t, expectedBody, rec.Body.String())
	assert.Equal(t, MIMEApplicationXMLCharsetUTF8, rec.Header().Get(HeaderContentType))
}
