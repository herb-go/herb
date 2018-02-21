package misc

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/herb-go/herb/middleware"
)

func TestElapsedtime(t *testing.T) {
	var successAction = func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			panic(err)
		}
	}
	var timespentMiddleware = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		time.Sleep(time.Second)
		next(w, r)
	}
	app := middleware.New()
	app.Use(ElapsedTime, timespentMiddleware).HandleFunc(successAction)
	server := httptest.NewServer(app)
	defer server.Close()
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	etstring := resp.Header.Get("Elapsed-Time")
	if len(etstring) < 4 || etstring[len(etstring)-3:] != " ns" {
		t.Error(etstring)
	}
	etstring = etstring[:len(etstring)-3]
	eti, err := strconv.Atoi(etstring)
	if err != nil {
		t.Fatal(err)
	}
	et := time.Duration(eti)
	if et < time.Second || et > 2*time.Second {
		t.Error(etstring)
	}
}
