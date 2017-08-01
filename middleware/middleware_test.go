package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testResults struct {
	data string
}

func (tr *testResults) newBeforeMiddleware(result string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		tr.data = tr.data + result
		next(w, r)
	}
}
func (tr *testResults) newAfterMiddleware(result string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		next(w, r)
		tr.data = tr.data + result
	}
}
func (tr *testResults) newHandleFunc(result string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tr.data = tr.data + result
	}
}
func TestNew(t *testing.T) {
	var ret testResults
	var app = New(ret.newBeforeMiddleware("a"), ret.newBeforeMiddleware("b"), ret.newBeforeMiddleware("c"))
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	var rr = httptest.NewRecorder()
	app.ServeHTTP(rr, r)
	if ret.data != "abc" {
		t.Errorf("New middleware order %s error", ret.data)
	}
}

func TestAfterNext(t *testing.T) {
	var ret testResults
	var app = New(ret.newAfterMiddleware("c"), ret.newAfterMiddleware("b"), ret.newAfterMiddleware("a"))
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	var rr = httptest.NewRecorder()
	app.ServeHTTP(rr, r)
	if ret.data != "abc" {
		t.Errorf("New middleware order %s error", ret.data)
	}
}

func TestUse(t *testing.T) {
	var ret testResults
	var app = New(ret.newBeforeMiddleware("a"), ret.newBeforeMiddleware("b"), ret.newBeforeMiddleware("c"))
	app.Use(ret.newAfterMiddleware("g"))
	app.Use(ret.newBeforeMiddleware("d"), ret.newBeforeMiddleware("e"))
	app.HandleFunc(ret.newHandleFunc("f"))
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	var rr = httptest.NewRecorder()
	app.ServeHTTP(rr, r)
	if ret.data != "abcdefg" {
		t.Errorf("New middleware order %s error", ret.data)
	}
}

func TestUseApp(t *testing.T) {
	var ret testResults
	var app = New(ret.newBeforeMiddleware("a"), ret.newBeforeMiddleware("b"), ret.newBeforeMiddleware("c"))
	var app2 = New(ret.newAfterMiddleware("f"))
	var app3 = New(ret.newBeforeMiddleware("d"))
	var app4 = New(ret.newBeforeMiddleware("e"))
	app.UseApp(app2)
	app.UseApp(app3, app4)
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	var rr = httptest.NewRecorder()
	app.ServeHTTP(rr, r)
	if ret.data != "abcdef" {
		t.Errorf("New middleware order %s error", ret.data)
	}
}
func TestHandlers(t *testing.T) {
	var ret testResults
	var app = New(ret.newBeforeMiddleware("a"), ret.newBeforeMiddleware("b"), ret.newBeforeMiddleware("c"))
	var length = len(app.Handlers())
	if length != 3 {
		t.Errorf("Middleware Handlers length %d error", length)
	}
	app.SetHandlers([]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){ret.newBeforeMiddleware("a")})
	var rr = httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	app.ServeHTTP(rr, r)
	if ret.data != "a" {
		t.Errorf("New middleware order %s error", ret.data)
	}
}

func TestServeMiddleware(t *testing.T) {
	var ret testResults
	var app = New(ret.newBeforeMiddleware("a"), ret.newBeforeMiddleware("b"), ret.newBeforeMiddleware("c"))
	var app2 = New(ret.newAfterMiddleware("f"))
	var app3 = New(ret.newBeforeMiddleware("d"))
	var app4 = New(ret.newBeforeMiddleware("e"))
	app.Use(app2.ServeMiddleware)
	app.Use(app3.ServeMiddleware, app4.ServeMiddleware)
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	var rr = httptest.NewRecorder()
	app.ServeHTTP(rr, r)
	if ret.data != "abcdef" {
		t.Errorf("New middleware order %s error", ret.data)
	}

}
func TestChain(t *testing.T) {
	var ret testResults
	var app = New(ret.newBeforeMiddleware("a"), ret.newBeforeMiddleware("b"), ret.newBeforeMiddleware("c"))
	var app2 = New(ret.newAfterMiddleware("e"))
	var app3 = New(ret.newBeforeMiddleware("d"))
	var app4 = New(ret.newBeforeMiddleware("f"))
	app.Chain(app2, app3)
	app2.UseApp(app3, app4)
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	var rr = httptest.NewRecorder()
	app.ServeHTTP(rr, r)
	if ret.data != "abcde" {
		t.Errorf("New middleware order %s error", ret.data)
	}
}
func TestFuncAppToMiddlewares(t *testing.T) {
	var ret testResults
	var app = New(ret.newBeforeMiddleware("a"), ret.newBeforeMiddleware("b"), ret.newBeforeMiddleware("c"))
	var app2 = New(AppToMiddlewares(app)...)
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	var rr = httptest.NewRecorder()
	app2.ServeHTTP(rr, r)
	if ret.data != "abc" {
		t.Errorf("New middleware order %s error", ret.data)
	}

}

func TestFuncUse(t *testing.T) {
	var ret testResults
	var app = New(ret.newBeforeMiddleware("a"), ret.newBeforeMiddleware("b"), ret.newBeforeMiddleware("c"))
	Use(app, ret.newAfterMiddleware("f"))
	Use(app, ret.newBeforeMiddleware("d"), ret.newBeforeMiddleware("e"))
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	var rr = httptest.NewRecorder()
	app.ServeHTTP(rr, r)
	if ret.data != "abcdef" {
		t.Errorf("New middleware order %s error", ret.data)
	}
}

func TestFuncChain(t *testing.T) {
	var ret testResults
	var app = New(ret.newBeforeMiddleware("a"), ret.newBeforeMiddleware("b"), ret.newBeforeMiddleware("c"))
	var app2 = New(ret.newAfterMiddleware("e"))
	var app3 = New(ret.newBeforeMiddleware("d"))
	var app4 = New(ret.newBeforeMiddleware("f"))
	Chain(app, app2, app3)
	app2.UseApp(app3, app4)
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	var rr = httptest.NewRecorder()
	app.ServeHTTP(rr, r)
	if ret.data != "abcde" {
		t.Errorf("New middleware order %s error", ret.data)
	}

}

func TestFuncWrapFunc(t *testing.T) {
	var ret testResults
	var app = New(ret.newBeforeMiddleware("a"), ret.newBeforeMiddleware("b"), ret.newBeforeMiddleware("c"))
	app.Use(ret.newAfterMiddleware("g"))
	app.Use(ret.newBeforeMiddleware("d"))
	app.Use(WrapFunc(ret.newHandleFunc("e")))
	app.HandleFunc(ret.newHandleFunc("f"))
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	var rr = httptest.NewRecorder()
	app.ServeHTTP(rr, r)
	if ret.data != "abcdefg" {
		t.Errorf("New middleware order %s error", ret.data)
	}

}
func TestFuncServeMiddleware(t *testing.T) {
	var ret testResults
	var app = New(ret.newBeforeMiddleware("a"), ret.newBeforeMiddleware("b"), ret.newBeforeMiddleware("c"))
	app.Use(ret.newAfterMiddleware("h"))
	app.Use(ret.newBeforeMiddleware("d"), ret.newBeforeMiddleware("e"))
	app.HandleFunc(ret.newHandleFunc("f"))
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	var rr = httptest.NewRecorder()
	ServeMiddleware(app, rr, r, ret.newHandleFunc("g"))
	if ret.data != "abcdefgh" {
		t.Errorf("New middleware order %s error", ret.data)
	}

}

func TestFuncServeHttp(t *testing.T) {
	var ret testResults
	var app = New(ret.newBeforeMiddleware("a"), ret.newBeforeMiddleware("b"), ret.newBeforeMiddleware("c"))
	app.Use(ret.newAfterMiddleware("g"))
	app.Use(ret.newBeforeMiddleware("d"), ret.newBeforeMiddleware("e"))
	app.HandleFunc(ret.newHandleFunc("f"))
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	var rr = httptest.NewRecorder()
	ServeHTTP(app, rr, r)
	if ret.data != "abcdefg" {
		t.Errorf("New middleware order %s error", ret.data)
	}

}
func TestFuncWrap(t *testing.T) {
	var ret testResults
	var app = New(ret.newBeforeMiddleware("a"), ret.newBeforeMiddleware("b"), ret.newBeforeMiddleware("c"))
	app.Use(ret.newAfterMiddleware("g"))
	app.Use(Wrap(http.HandlerFunc((ret.newHandleFunc("d")))))
	app.Use(ret.newBeforeMiddleware("e"))
	app.HandleFunc(ret.newHandleFunc("f"))
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	var rr = httptest.NewRecorder()
	ServeHTTP(app, rr, r)
	if ret.data != "abcdefg" {
		t.Errorf("New middleware order %s error", ret.data)
	}

}
