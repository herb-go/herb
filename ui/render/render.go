package render

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

//const from https://github.com/unrolled/render/blob/v1/render.go
const (
	// ContentBinary header value for binary data.
	ContentBinary = "application/octet-stream"
	// ContentHTML header value for HTML data.
	ContentHTML = "text/html"
	// ContentJSON header value for JSON data.
	ContentJSON = "application/json"
	// ContentJSONP header value for JSONP data.
	ContentJSONP = "application/javascript"
	// ContentLength header constant.
	ContentLength = "Content-Length"
	// ContentText header value for Text data.
	ContentText = "text/plain"
	// ContentType header constant.
	ContentType = "Content-Type"
	// ContentXHTML header value for XHTML data.
	ContentXHTML = "application/xhtml+xml"
	// ContentXML header value for XML data.
	ContentXML = "text/xml"
	// Default character encoding.
	defaultCharset = "UTF-8"
)

//ErrTooManyViewFiles error rasied when too many view files is given.
//raised by render engine.
var ErrTooManyViewFiles = errors.New("error too many view files")

//WriteJSON write json data to response.
//Return bytes length wrote and any error if raised.
func WriteJSON(w http.ResponseWriter, data []byte, status int) (int, error) {
	w.Header().Set(ContentType, ContentJSON)
	w.WriteHeader(status)
	return w.Write(data)
}

//MustWriteJSON write json data to response.
//Return bytes length wrote.
//Panic if any error raised.
func MustWriteJSON(w http.ResponseWriter, data []byte, status int) int {
	result, err := WriteJSON(w, data, status)
	if err != nil {
		panic(err)
	}
	return result
}

//WriteHTML write html data to response.
//Return bytes length wrote and any error if raised.
func WriteHTML(w http.ResponseWriter, data []byte, status int) (int, error) {
	w.Header().Set(ContentType, ContentHTML)
	w.WriteHeader(status)
	return w.Write(data)
}

//MustWriteHTML write html data to response.
//Return bytes length wrote.
//Panic if any error raised.
func MustWriteHTML(w http.ResponseWriter, data []byte, status int) int {
	result, err := WriteHTML(w, data, status)
	if err != nil {
		panic(err)
	}
	return result
}

//HTMLFile write content of given file to response as html.
//Return bytes length wrote and any error if raised.
func HTMLFile(w http.ResponseWriter, path string, status int) (int, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return WriteHTML(w, bytes, status)
}

//MustHTMLFile write content of given file to response as html.
//Return bytes length wrote.
//Panic if any error raised.
func MustHTMLFile(w http.ResponseWriter, path string, status int) int {
	result, err := HTMLFile(w, path, status)
	if err != nil {
		panic(err)
	}
	return result
}

//JSON marshal data as json and write to response
//Return bytes length wrote and any error if raised.
func JSON(w http.ResponseWriter, data interface{}, status int) (int, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return WriteJSON(w, bytes, status)
}

//MustJSON marshal data as json and write to response
//Return bytes length wrote.
//Panic if any error raised.
func MustJSON(w http.ResponseWriter, data interface{}, status int) int {
	result, err := JSON(w, data, status)
	if err != nil {
		panic(err)
	}
	return result
}

//Error write a http error to response
//Return bytes length wrote.
//Panic if any error raised.
func Error(w http.ResponseWriter, status int) (int, error) {
	w.Header().Set(ContentType, ContentText)
	w.WriteHeader(status)
	return w.Write([]byte(http.StatusText(status)))
}

//MustError write a http error to response
//Return bytes length wrote.
//Panic if any error raised.
func MustError(w http.ResponseWriter, status int) int {
	result, err := Error(w, status)
	if err != nil {
		panic(err)
	}
	return result
}

//JSONHandler create a json handler with give data and status code
func JSONHandler(data interface{}, status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MustJSON(w, data, status)
	})
}

//HTMLHandler create a html handler with give data and status code
func HTMLHandler(data string, status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MustWriteHTML(w, []byte(data), status)
	})
}
