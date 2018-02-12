package httprouter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware/router"
)

func testAction(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func TestRouter(t *testing.T) {
	router := New()
	router.GET("/get").HandleFunc(testAction)
	router.POST("/post").HandleFunc(testAction)
	router.PUT("/put").HandleFunc(testAction)
	router.DELETE("/delete").HandleFunc(testAction)
	router.HEAD("/head").HandleFunc(testAction)
	router.OPTIONS("/options").HandleFunc(testAction)
	router.PATCH("/patch").HandleFunc(testAction)
	router.ALL("/all").HandleFunc(testAction)
	server := httptest.NewServer(router)
	defer server.Close()
	var tests = []string{"/get", "/post", "/put", "/delete", "/head", "/options", "/patch", "/all"}
	var methods = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"}
	var result = map[string]map[string]int{}
	for _, v := range tests {
		result[v] = map[string]int{}
		for _, method := range methods {
			req, err := http.NewRequest(method, server.URL+v, nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			result[v][method] = resp.StatusCode
		}
	}
	var testSuites = map[string]string{
		"/get":    "GET",
		"/post":   "POST",
		"/put":    "PUT",
		"/delete": "DELETE",
		"/head":   "HEAD",
	}
	for url := range testSuites {
		for method := range result[url] {
			if method == testSuites[url] {
				if result[url][method] != 200 {
					t.Error(url, method, result[url][method])
				}
			} else if method == "OPTIONS" {
				if result[url][method] != 200 {
					t.Error(url, method, result[url][method])
				}
			} else if result[url][method] != 405 {
				t.Error(url, method, result[url][method])
			}
		}
	}
	var url = "/options"
	for method := range result[url] {
		if method == "OPTIONS" {
			if result[url][method] != 200 {
				t.Error(url, method, result[url][method])
			}
		} else if result[url][method] != 404 {
			t.Error(url, method, result[url][method])
		}
	}
	url = "/all"
	for method := range result[url] {
		if result[url][method] != 200 {
			t.Error(url, method, result[url][method])
		}
	}

}

func TestStripPrefix(t *testing.T) {
	router := New()
	subrouter := New()
	subrouter.GET("/action").HandleFunc(testAction)
	router.StripPrefix("/test").Handle(subrouter)
	server := httptest.NewServer(router)
	defer server.Close()
	req, err := http.NewRequest("GET", server.URL+"/test/action", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
}

func TestParam(t *testing.T) {
	router := New()
	router.GET("/test/:id").HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(GetParams(r).Get("id")))
		if err != nil {
			panic(err)
		}
	})
	server := httptest.NewServer(router)
	defer server.Close()
	req, err := http.NewRequest("GET", server.URL+"/test/action", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "action" {
		t.Error(string(content))
	}
}

func TestRouterParam(t *testing.T) {
	r := New()
	r.GET("/test/:id").HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(router.GetParams(r).Get("id")))
		if err != nil {
			panic(err)
		}
	})
	server := httptest.NewServer(r)
	defer server.Close()
	req, err := http.NewRequest("GET", server.URL+"/test/action", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "action" {
		t.Error(string(content))
	}
}
