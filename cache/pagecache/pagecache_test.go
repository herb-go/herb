package pagecache

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	_ "github.com/herb-go/herb/cache/drivers/syncmapcache"

	"github.com/herb-go/herb/cache"
)

func newTestCache(ttl int64) *cache.Cache {
	config := &cache.ConfigJSON{}
	config.Set("Size", 10000000)
	c := cache.New()
	oc := &cache.OptionConfigMap{
		Driver:    "syncmapcache",
		TTL:       ttl * int64(time.Second),
		Config:    nil,
		Marshaler: "json",
	}
	err := c.Init(oc)
	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	return c

}

var content int

func testAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("ts", strconv.FormatInt(time.Now().UnixNano(), 10))
	w.WriteHeader(content)
	w.Write([]byte(r.Header.Get("test") + strconv.Itoa(content)))
}
func emptyKeyGenerator(r *http.Request) string {
	return ""
}
func keyGenerator(r *http.Request) string {
	return r.Header.Get("test")
}
func fieldGenerator(c *cache.Cache) func(r *http.Request) *cache.Field {
	return func(r *http.Request) *cache.Field {
		var test = r.Header.Get("test")
		return c.Field(test)
	}
}
func nilFieldGenerator(r *http.Request) *cache.Field {
	return nil
}
func TestPageCacheField(t *testing.T) {
	content = 200
	mux := http.NewServeMux()
	fg := fieldGenerator(newTestCache(3600))
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("rawts", strconv.FormatInt(time.Now().UnixNano(), 10))
		FieldMiddleware(fg, 3600*time.Second, nil)(w, r, testAction)
	})
	mux.HandleFunc("/raw", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("rawts", strconv.FormatInt(time.Now().UnixNano(), 10))
		testAction(w, r)
	})
	mux.HandleFunc("/nil", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("rawts", strconv.FormatInt(time.Now().UnixNano(), 10))
		FieldMiddleware(nilFieldGenerator, 3600*time.Second, nil)(w, r, testAction)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	req, err := http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Set("test", "test1")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	ts := resp.Header.Get("ts")
	rawts := resp.Header.Get("rawts")
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test1"+"200" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/raw", nil)
	req.Header.Set("test", "test1")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test1"+"200" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/nil", nil)
	req.Header.Set("test", "test1")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test1"+"200" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	content = 404
	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Set("test", "test1")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test1"+"200" {
		t.Error(string(bs))
	}
	if resp.Header.Get("ts") != ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/raw", nil)
	req.Header.Set("test", "test1")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test1"+"404" {
		t.Error(string(bs))
	}
	if resp.Header.Get("ts") == ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/nil", nil)
	req.Header.Set("test", "test1")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test1"+"404" {
		t.Error(string(bs))
	}
	if resp.Header.Get("ts") == ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Set("test", "test2")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test2"+"404" {
		t.Error(string(bs))
	}
	if resp.Header.Get("ts") == ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}
	content = 500
	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Set("test", "test3")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test3"+"500" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 500 {
		t.Error(resp.StatusCode)
	}
	content = 403
	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Set("test", "test3")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test3"+"403" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 403 {
		t.Error(resp.StatusCode)
	}
}

func TestPageCache(t *testing.T) {
	var PageCache = New(newTestCache(3600))
	content = 200
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("rawts", strconv.FormatInt(time.Now().UnixNano(), 10))
		PageCache.Middleware(keyGenerator, 0)(w, r, testAction)
	})
	mux.HandleFunc("/raw", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("rawts", strconv.FormatInt(time.Now().UnixNano(), 10))
		testAction(w, r)
	})
	mux.HandleFunc("/empty", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("rawts", strconv.FormatInt(time.Now().UnixNano(), 10))
		PageCache.Middleware(emptyKeyGenerator, 0)(w, r, testAction)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	req, err := http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Set("test", "test1")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	ts := resp.Header.Get("ts")
	rawts := resp.Header.Get("rawts")
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test1"+"200" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/raw", nil)
	req.Header.Set("test", "test1")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test1"+"200" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/empty", nil)
	req.Header.Set("test", "test1")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test1"+"200" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	content = 404
	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Set("test", "test1")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test1"+"200" {
		t.Error(string(bs))
	}
	if resp.Header.Get("ts") != ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/raw", nil)
	req.Header.Set("test", "test1")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test1"+"404" {
		t.Error(string(bs))
	}
	if resp.Header.Get("ts") == ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/empty", nil)
	req.Header.Set("test", "test1")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test1"+"404" {
		t.Error(string(bs))
	}
	if resp.Header.Get("ts") == ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Set("test", "test2")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test2"+"404" {
		t.Error(string(bs))
	}
	if resp.Header.Get("ts") == ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}
	content = 500
	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Set("test", "test3")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test3"+"500" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 500 {
		t.Error(resp.StatusCode)
	}
	content = 403
	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Set("test", "test3")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "test3"+"403" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 403 {
		t.Error(resp.StatusCode)
	}
}
