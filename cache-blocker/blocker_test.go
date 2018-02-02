package blocker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/herb-go/herb/cache"
	_ "github.com/herb-go/herb/cache/drivers/freecache"
)

func newTestCache(ttl int64) *cache.Cache {
	config := json.RawMessage("{\"Size\": 10000000}")
	c := cache.New()
	err := c.Open("freecache", config, ttl)
	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	return c
}
func testIdentifier(r *http.Request) (string, error) {
	return r.Header.Get("name"), nil
}
func TestBlock(t *testing.T) {
	var rep *http.Response
	var err error
	blocker := New(newTestCache(1*3600), testIdentifier)
	blocker.Block(0, 20, 1*time.Hour)
	blocker.Block(404, 5, 1*time.Hour)
	blocker.Block(403, 5, 1*time.Hour)
	mux := http.NewServeMux()
	mux.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("/403", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(403), 403)
	})
	mux.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(404), 404)
	})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blocker.ServeMiddleware(w, r, mux.ServeHTTP)
	}))
	defer server.Close()
	req, err := http.NewRequest("get", server.URL+"/403", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("name", "test1")
	for i := 0; i < 5; i++ {
		rep, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		rep.Body.Close()
		if rep.StatusCode != 403 {
			t.Error(rep.StatusCode)
		}
		time.Sleep(10 * time.Millisecond)
	}
	rep, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	rep.Body.Close()
	if rep.StatusCode != 429 {
		t.Error(rep.StatusCode)
	}
	time.Sleep(10 * time.Millisecond)
	req, err = http.NewRequest("get", server.URL+"/403", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("name", "test2")
	rep, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	rep.Body.Close()
	if rep.StatusCode != 403 {
		t.Error(rep.StatusCode)
	}
	time.Sleep(10 * time.Millisecond)
	blocker.StatusBlocked = 400
	req, err = http.NewRequest("get", server.URL+"/404", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("name", "test2")
	for i := 0; i < 5; i++ {
		rep, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		rep.Body.Close()
		if rep.StatusCode != 404 {
			t.Error(rep.StatusCode)
		}
		time.Sleep(10 * time.Millisecond)
	}
	rep, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	rep.Body.Close()
	if rep.StatusCode != 400 {
		t.Error(rep.StatusCode)
	}
	time.Sleep(10 * time.Millisecond)
	blocker.StatusBlocked = defaultBlockedStatus
	req, err = http.NewRequest("get", server.URL+"/403", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("name", "test3")
	for i := 0; i < 4; i++ {
		rep, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		rep.Body.Close()
		if rep.StatusCode != 403 {
			t.Error(rep.StatusCode)
		}
		time.Sleep(10 * time.Millisecond)
	}

	req, err = http.NewRequest("get", server.URL+"/404", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("name", "test3")
	for i := 0; i < 4; i++ {
		rep, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		rep.Body.Close()
		if rep.StatusCode != 404 {
			t.Error(rep.StatusCode)
		}
		time.Sleep(10 * time.Millisecond)
	}

	req, err = http.NewRequest("get", server.URL+"/200", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("name", "test3")
	for i := 0; i < 12; i++ {
		rep, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		rep.Body.Close()
		if rep.StatusCode != 200 {
			t.Error(rep.StatusCode)
		}
		time.Sleep(10 * time.Millisecond)
	}
	rep, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	rep.Body.Close()
	if rep.StatusCode != 429 {
		t.Error(rep.StatusCode)
	}
}

func TestIPIdentifier(t *testing.T) {
	var rep *http.Response
	var err error
	blocker := New(newTestCache(1*3600), IPIdentifier)
	blocker.Block(0, 20, 1*time.Hour)
	blocker.Block(404, 5, 1*time.Hour)
	blocker.Block(403, 5, 1*time.Hour)
	mux := http.NewServeMux()
	mux.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("/403", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(403), 403)
	})
	mux.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(404), 404)
	})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blocker.ServeMiddleware(w, r, mux.ServeHTTP)
	}))
	defer server.Close()
	req, err := http.NewRequest("get", server.URL+"/403", nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 5; i++ {
		rep, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		rep.Body.Close()
		if rep.StatusCode != 403 {
			t.Error(rep.StatusCode)
		}
		time.Sleep(10 * time.Millisecond)
	}
	rep, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	rep.Body.Close()
	if rep.StatusCode != 429 {
		t.Error(rep.StatusCode)
	}
	time.Sleep(10 * time.Millisecond)
}
