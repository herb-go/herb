package render

import (
	"encoding/json"
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

func WriteJSON(w http.ResponseWriter, data []byte, status int) (int, error) {
	w.Header().Set(ContentType, ContentJSON)
	w.WriteHeader(status)
	return w.Write(data)
}
func MustWriteJSON(w http.ResponseWriter, data []byte, status int) int {
	result, err := WriteJSON(w, data, status)
	if err != nil {
		panic(err)
	}
	return result
}
func WriteHTML(w http.ResponseWriter, data []byte, status int) (int, error) {
	w.Header().Set(ContentType, ContentHTML)
	w.WriteHeader(status)
	return w.Write(data)
}

func MustWriteHTML(w http.ResponseWriter, data []byte, status int) int {
	w.Header().Set(ContentType, ContentHTML)
	w.WriteHeader(status)
	result, err := w.Write(data)
	if err != nil {
		panic(err)
	}
	return result
}

func JSON(w http.ResponseWriter, data interface{}, status int) (int, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return WriteJSON(w, bytes, status)
}
func MustJSON(w http.ResponseWriter, data interface{}, status int) int {
	result, err := JSON(w, data, status)
	if err != nil {
		panic(err)
	}
	return result
}
func Error(w http.ResponseWriter, status int) (int, error) {
	w.Header().Set(ContentType, ContentText)
	w.WriteHeader(status)
	return w.Write([]byte(http.StatusText(status)))
}

func MustError(w http.ResponseWriter, status int) int {
	result, err := Error(w, status)
	if err != nil {
		panic(err)
	}
	return result
}
