package session

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/herb-go/herb/cache"
	_ "github.com/herb-go/herb/cache/drivers/syncmapcache"
	_ "github.com/herb-go/herb/cache/marshalers/msgpackmarshaler"
)

func getStore(ttl time.Duration) *Store {
	c := cache.New()
	oc := cache.NewOptionConfig()
	oc.Driver = "syncmapcache"
	oc.TTL = int64(ttl) * int64(time.Second)
	oc.Config = nil
	oc.Marshaler = "json"
	err := c.Init(oc)

	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	s := MustCacheStore(c, ttl)
	m, err := cache.NewMarshaler("json")
	if err != nil {
		panic(err)
	}
	s.Marshaler = m
	return s
}

func getTimeoutStore(ttl time.Duration, UpdateActiveInterval time.Duration) *Store {
	c := cache.New()
	oc := cache.NewOptionConfig()
	oc.Driver = "syncmapcache"
	oc.TTL = int64(ttl) * int64(time.Second)
	oc.Config = nil
	oc.Marshaler = "json"
	err := c.Init(oc)
	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	s := MustCacheStore(c, ttl)
	s.UpdateActiveInterval = UpdateActiveInterval
	m, err := cache.NewMarshaler("json")
	if err != nil {
		panic(err)
	}
	s.Marshaler = m
	return s
}
func TestField(t *testing.T) {
	var err error
	s := getStore(time.Hour)
	defer s.Close()
	model := "123456"
	model2 := "abcde"
	var result string
	testKey := "testkey"
	testOwner := "testowner"
	testHeaderName := "token"
	testuid := "testuid"

	field := s.Field(testKey)
	token, err := s.GenerateToken(testOwner)
	if err != nil {
		panic(err)
	}
	ts := s.GetSession(token)
	err = field.SaveTo(ts, model)
	if err != nil {
		panic(err)
	}
	err = ts.Save()
	if err != nil {
		panic(err)
	}
	err = ts.Get(testKey, &result)
	if err != nil {
		panic(err)
	}
	if result != model {
		t.Errorf("Reuslt error %s", result)
	}
	result = ""
	err = field.LoadFrom(ts, &result)
	if err != nil {
		panic(err)
	}
	if result != model {
		t.Errorf("Reuslt error %s", result)
	}
	ActionTestField := func(w http.ResponseWriter, r *http.Request) {
		result = ""
		err := field.Get(r, &result)
		if err != nil {
			panic(err)
		}
		if result != model {
			t.Errorf("Reuslt error %s", result)
		}
		err = field.Flush(r)
		if err != nil {
			panic(err)
		}
		err = field.Get(r, &result)
		if err != ErrDataNotFound {
			panic(err)
		}
		session, err := field.GetSession(r)
		if session.token != ts.token {
			t.Fatal(session)
		}
		if err != nil {
			panic(err)
		}
		w.Write([]byte("ok"))
	}
	ActionLogin := func(w http.ResponseWriter, r *http.Request) {
		uid := r.Header.Get("uid")
		err = field.Login(w, r, uid)
		if err != nil {
			panic(err)
		}
		s, err := field.GetSession(r)
		if err != nil {
			panic(err)
		}
		ts.token = s.token
		w.Write([]byte(uid))
	}
	ActionUID := func(w http.ResponseWriter, r *http.Request) {
		uid, err := field.IdentifyRequest(r)
		if err != nil {
			panic(err)
		}
		w.Write([]byte(uid))
	}
	ActionLogout := func(w http.ResponseWriter, r *http.Request) {
		err := field.Logout(w, r)
		if err != nil {
			panic(err)
		}
		w.Write([]byte("ok"))
	}
	ActionGetField := func(w http.ResponseWriter, r *http.Request) {
		result = ""
		err := field.Get(r, &result)
		if err != nil {
			panic(err)
		}
		if result != model {
			t.Errorf("Reuslt error %s", result)
		}
	}
	ActionSetField := func(w http.ResponseWriter, r *http.Request) {
		err := field.Set(r, model2)
		if err != nil {
			panic(err)
		}
	}
	var mux = http.NewServeMux()
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) { s.HeaderMiddleware(testHeaderName)(w, r, ActionGetField) })
	mux.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) { s.HeaderMiddleware(testHeaderName)(w, r, ActionSetField) })
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		s.HeaderMiddleware(testHeaderName)(w, r, ActionTestField)
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) { s.HeaderMiddleware(testHeaderName)(w, r, ActionLogin) })
	mux.HandleFunc("/uid", func(w http.ResponseWriter, r *http.Request) { s.HeaderMiddleware(testHeaderName)(w, r, ActionUID) })
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) { s.HeaderMiddleware(testHeaderName)(w, r, ActionLogout) })
	hs := httptest.NewServer(mux)
	defer hs.Close()
	c := &http.Client{}
	GetRequest, err := http.NewRequest("POST", hs.URL+"/get", nil)
	GetRequest.Header.Set(testHeaderName, ts.MustToken())
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Do(GetRequest)
	if err != nil {
		t.Fatal(err)
	}
	TestRequest, err := http.NewRequest("POST", hs.URL+"/test", nil)
	TestRequest.Header.Set(testHeaderName, ts.MustToken())
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Do(TestRequest)
	if err != nil {
		t.Fatal(err)
	}

	SetRequest, err := http.NewRequest("POST", hs.URL+"/set", nil)
	SetRequest.Header.Set(testHeaderName, ts.MustToken())
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Do(SetRequest)
	if err != nil {
		t.Fatal(err)
	}
	result = ""
	ts = s.GetSession(ts.MustToken())
	err = field.LoadFrom(ts, &result)
	if err != nil {
		panic(err)
	}
	if result != model2 {
		t.Errorf("Reuslt error %s", result)
	}
	LoginRequest, err := http.NewRequest("POST", hs.URL+"/login", nil)
	LoginRequest.Header.Set("uid", testuid)
	LoginRequest.Header.Set(testHeaderName, ts.MustToken())
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Do(LoginRequest)
	if err != nil {
		t.Fatal(err)
	}

	UIDRequest, err := http.NewRequest("POST", hs.URL+"/uid", nil)
	UIDRequest.Header.Set(testHeaderName, ts.MustToken())
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Do(UIDRequest)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(body) != testuid {
		t.Fatal(string(body))
	}

	LogoutRequest, err := http.NewRequest("POST", hs.URL+"/logout", nil)
	LogoutRequest.Header.Set(testHeaderName, ts.MustToken())
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Do(LogoutRequest)
	if err != nil {
		t.Fatal(err)
	}

	UIDRequest, err = http.NewRequest("POST", hs.URL+"/uid", nil)
	UIDRequest.Header.Set(testHeaderName, ts.MustToken())
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(UIDRequest)
	if err != nil {
		t.Fatal(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(body) != "" {
		t.Fatal(string(body))
	}
}
func TestTD(t *testing.T) {
	var err error
	s := getStore(time.Hour)
	defer s.Close()
	model := "123456"
	var result string
	testKey := "testkey"
	type modelStruct struct {
		Data string
	}
	structModel := modelStruct{
		Data: "test",
	}
	var resutStruct = modelStruct{}
	var testStructKey = "teststructkey"
	var modelInt = 123456
	var resultInt int
	var testIntKey = "testintkey"
	var modelBytes = []byte("testbytes")
	var resultBytes []byte
	var testBytesKey = "testbyteskey"
	var modelMap = map[string]string{
		"test": "test",
	}
	var resultMap map[string]string
	var testMapKey = "testmapkey"
	testOwner := "testowner"

	td, err := s.RegenerateToken(testOwner)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Set(testKey, model)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Save()
	if err != nil {
		t.Fatal(err)
	}
	result = ""
	err = td.Get(testKey, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result != model {
		t.Errorf("td  error %s", result)
	}
	result = ""
	td = s.GetSession(td.MustToken())
	err = s.LoadSession(td)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Get(testKey, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result != model {
		t.Errorf("td LoadSession error")
	}

	td, err = s.RegenerateToken(testOwner)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Set(testStructKey, structModel)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Save()
	if err != nil {
		t.Fatal(err)
	}
	resutStruct = modelStruct{}
	err = td.Get(testStructKey, &resutStruct)
	if err != nil {
		t.Fatal(err)
	}
	if resutStruct != structModel {
		t.Errorf("td  error %s", resutStruct)
	}
	resutStruct = modelStruct{}
	td = s.GetSession(td.MustToken())
	err = s.LoadSession(td)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Get(testStructKey, &resutStruct)
	if err != nil {
		t.Fatal(err)
	}
	if resutStruct != structModel {
		t.Errorf("td LoadSession error")
	}

	td, err = s.RegenerateToken(testOwner)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Set(testIntKey, modelInt)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Save()
	if err != nil {
		t.Fatal(err)
	}
	resultInt = 0
	err = td.Get(testIntKey, &resultInt)
	if err != nil {
		t.Fatal(err)
	}
	if resultInt != modelInt {
		t.Errorf("td  error %d", resultInt)
	}
	resultInt = 0
	td = s.GetSession(td.MustToken())
	err = s.LoadSession(td)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Get(testIntKey, &resultInt)
	if err != nil {
		t.Fatal(err)
	}
	if resultInt != modelInt {
		t.Errorf("td LoadSession error")
	}

	td, err = s.RegenerateToken(testOwner)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Set(testBytesKey, modelBytes)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Save()
	if err != nil {
		t.Fatal(err)
	}
	resultBytes = []byte{}
	err = td.Get(testBytesKey, &resultBytes)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(resultBytes, modelBytes) != 0 {
		t.Errorf("td  error %s", resultBytes)
	}
	resultBytes = []byte{}
	td = s.GetSession(td.MustToken())
	err = s.LoadSession(td)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Get(testBytesKey, &resultBytes)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(resultBytes, modelBytes) != 0 {
		t.Errorf("td  error %s", resultBytes)
	}

	td, err = s.RegenerateToken(testOwner)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Set(testMapKey, modelMap)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Save()
	if err != nil {
		t.Fatal(err)
	}
	resultMap = map[string]string{}
	err = td.Get(testMapKey, &resultMap)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(resultMap, modelMap) {
		t.Errorf("td  error %s", resultMap)
	}
	resultMap = map[string]string{}
	td = s.GetSession(td.MustToken())
	err = s.LoadSession(td)
	if err != nil {
		t.Fatal(err)
	}
	err = td.Get(testMapKey, &resultMap)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(resultMap, modelMap) {
		t.Errorf("td  error %s", resultMap)
	}

}

func TestTimeout(t *testing.T) {
	sforever := getTimeoutStore(time.Hour, -1)
	s3second := getTimeoutStore(3*time.Second, -1)
	s3secondwithAutoRefresh := getTimeoutStore(3*time.Second, 1*time.Second)
	testOwner := "testowner"
	model := "123456"
	var result string
	testKey := "testkey"
	tdForeverKey, err := sforever.GenerateToken(testOwner)
	if err != nil {
		panic(err)
	}
	tdForever, err := sforever.GenerateSession(tdForeverKey)
	if err != nil {
		panic(err)
	}
	err = tdForever.Set(testKey, model)
	if err != nil {
		panic(err)
	}
	td3secondKey, err := s3second.GenerateToken(testOwner)
	if err != nil {
		panic(err)
	}
	td3second, err := s3second.GenerateSession(td3secondKey)
	if err != nil {
		panic(err)
	}
	err = td3second.Set(testKey, model)
	if err != nil {
		panic(err)
	}
	td3secondwithAutoRefreshKey, err := s3secondwithAutoRefresh.GenerateToken(testOwner)
	if err != nil {
		panic(err)
	}
	td3secondwithAutoRefresh, err := s3secondwithAutoRefresh.GenerateSession(td3secondwithAutoRefreshKey)
	if err != nil {
		panic(err)
	}
	err = td3secondwithAutoRefresh.Set(testKey, model)
	if err != nil {
		panic(err)
	}
	tdForever.Save()
	td3second.Save()
	td3secondwithAutoRefresh.Save()
	time.Sleep(2 * time.Second)
	tdForever = sforever.GetSession(tdForever.MustToken())

	td3second = s3second.GetSession(td3second.MustToken())
	td3secondwithAutoRefresh = s3secondwithAutoRefresh.GetSession(td3secondwithAutoRefresh.MustToken())
	result = ""
	err = tdForever.Get(testKey, &result)
	if result != model {
		t.Errorf("Timeout error %s", result)
	}
	result = ""
	err = td3second.Get(testKey, &result)
	if result != model {
		t.Errorf("Timeout error %s", result)
	}
	result = ""
	err = td3secondwithAutoRefresh.Get(testKey, &result)
	if result != model {
		t.Errorf("Timeout error %s", result)
	}
	tdForever.Save()
	td3second.Save()
	td3secondwithAutoRefresh.Save()
	time.Sleep(2 * time.Second)
	tdForever = sforever.GetSession(tdForever.MustToken())
	td3second = s3second.GetSession(td3second.MustToken())
	td3secondwithAutoRefresh = s3secondwithAutoRefresh.GetSession(td3secondwithAutoRefresh.MustToken())

	result = ""
	err = tdForever.Get(testKey, &result)
	if result != model {
		t.Errorf("Timeout error %s", result)
	}
	result = ""
	err = td3second.Get(testKey, &result)
	if err != ErrDataNotFound {
		t.Errorf("Timeout error %s", err)
	}
	result = ""
	err = td3secondwithAutoRefresh.Get(testKey, &result)
	if result != model {
		t.Errorf("Timeout error %s", result)
	}
	tdForever.Save()
	td3second.Save()
	td3secondwithAutoRefresh.Save()
	time.Sleep(4 * time.Second)
	tdForever = sforever.GetSession(tdForever.MustToken())
	td3second = s3second.GetSession(td3second.MustToken())

	td3secondwithAutoRefresh = s3secondwithAutoRefresh.GetSession(td3secondwithAutoRefresh.MustToken())

	result = ""
	err = tdForever.Get(testKey, &result)
	if result != model {
		t.Errorf("Timeout error %s", result)
	}
	result = ""
	err = td3second.Get(testKey, &result)
	if err != ErrDataNotFound {
		t.Errorf("Timeout error %s", err)
	}
	result = ""
	err = td3secondwithAutoRefresh.Get(testKey, &result)
	if err != ErrDataNotFound {
		t.Errorf("Timeout error %s", err)
	}
}
func TestSessionMarshal(t *testing.T) {
	var err error
	testOwner := "testowner"
	model := "123456"
	var result string
	testKey := "testkey"
	testKey2 := "testkey2"
	testToken := "testtoken"
	s := getStore(time.Hour)
	defer s.Close()
	td, err := s.GenerateSession(testOwner)
	if err != nil {
		panic(err)
	}
	err = td.Set(testKey, model)
	if err != nil {
		panic(err)
	}
	bytes, err := td.Marshal()
	if err != nil {
		panic(err)
	}
	td2 := NewSession(testToken, s)
	err = td2.Unmarshal(testKey2, bytes)
	if err != nil {
		panic(err)
	}
	err = td2.Get(testKey, &result)
	if err != nil {
		panic(err)
	}
	if result != model {
		t.Errorf("Session Unmarshal err %s", result)
	}

}

func TestRequest(t *testing.T) {
	var err error
	s := getStore(time.Hour)
	defer s.Close()
	model := "123456"
	modelAfterSet := "set"
	var result string
	testKey := "testkey"
	testOwner := "testowner"
	testHeaderName := "token"
	var token string
	var mux = http.NewServeMux()
	actionTest := func(w http.ResponseWriter, r *http.Request) {
		ts, err := s.GetRequestSession(r)
		if err != nil {
			t.Fatal(err)
		}
		s.Get(r, testKey, &result)
		if result != model {
			t.Errorf("Field get error %s", result)
		}
		result = ""
		ts, err = s.GenerateSession(ts.MustToken())
		if err != nil {
			t.Fatal(err)
		}
		err = ts.Get(testKey, &result)
		if err != nil {
			t.Fatal(err)
		}
		if result != model {
			t.Errorf("Field get error result %s", result)
		}
		ex, err := s.ExpiredAt(r)
		if err != nil {
			t.Fatal(err)
		}
		if ex <= 0 {
			t.Errorf("Field ExpiredAt error %d", ex)
		}
		err = s.Set(r, testKey, modelAfterSet)
		if err != nil {
			t.Fatal(err)
		}
		result = ""
		err = s.Get(r, testKey, &result)
		if err != nil {
			t.Fatal(err)
		}
		if result != modelAfterSet {
			t.Errorf("field.Set error %s", result)
		}
		w.Write([]byte("ok"))
	}
	actionHeaderTest := func(w http.ResponseWriter, r *http.Request) {
		s.HeaderMiddleware(testHeaderName)(w, r, actionTest)
	}
	actionCookieTest := func(w http.ResponseWriter, r *http.Request) {
		s.CookieMiddleware()(w, r, actionTest)
	}
	actionLogin := func(w http.ResponseWriter, r *http.Request) {
		ts, err := s.GetRequestSession(r)
		if err != nil {
			panic(err)
		}
		err = ts.RegenerateToken(testOwner)
		if err != nil {
			panic(err)
		}
		ts.Set(testKey, model)
		w.Write([]byte(ts.MustToken()))
	}
	actionHeaderLogin := func(w http.ResponseWriter, r *http.Request) {
		s.HeaderMiddleware(testHeaderName)(w, r, actionLogin)
	}
	actionCookieLogin := func(w http.ResponseWriter, r *http.Request) {
		s.CookieMiddleware()(w, r, actionLogin)
	}
	actionLogout := func(w http.ResponseWriter, r *http.Request) {
		s.DestoryMiddleware()(w, r, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	}
	actionCookieLogout := func(w http.ResponseWriter, r *http.Request) {
		s.CookieMiddleware()(w, r, actionLogout)
	}
	actionHeaderLogout := func(w http.ResponseWriter, r *http.Request) {
		s.HeaderMiddleware(testHeaderName)(w, r, actionLogout)
	}
	actionLoginStatus := func(w http.ResponseWriter, r *http.Request) {
		t, err := s.GetRequestSession(r)
		if err != ErrDataNotFound && err != nil {
			panic(err)
		}
		if err == ErrDataNotFound || t.token == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return

		}
		result = ""
		err = s.Get(r, testKey, &result)
		if err != ErrDataNotFound && err != nil {
			panic(err)
		}
		if err == ErrDataNotFound || result == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		w.Write([]byte("ok"))
	}
	actionCookieLoginStatus := func(w http.ResponseWriter, r *http.Request) {
		s.CookieMiddleware()(w, r, actionLoginStatus)
	}
	actionHeaderLoginStatus := func(w http.ResponseWriter, r *http.Request) {
		s.HeaderMiddleware(testHeaderName)(w, r, actionLoginStatus)
	}
	mux.HandleFunc("/login", actionHeaderLogin)
	mux.HandleFunc("/cookie/login", actionCookieLogin)
	mux.HandleFunc("/test", actionHeaderTest)
	mux.HandleFunc("/cookie/test", actionCookieTest)
	mux.HandleFunc("/logout", actionHeaderLogout)
	mux.HandleFunc("/cookie/logout", actionCookieLogout)
	mux.HandleFunc("/loginstatus", actionHeaderLoginStatus)
	mux.HandleFunc("/cookie/loginstatus", actionCookieLoginStatus)
	hs := httptest.NewServer(mux)
	defer hs.Close()
	c := &http.Client{}
	LoginStatusRequest, err := http.NewRequest("POST", hs.URL+"/loginstatus", nil)
	if err != nil {
		t.Fatal(err)
	}
	rep, err := c.Do(LoginStatusRequest)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != http.StatusUnauthorized {
		t.Errorf("Status code error %d", rep.StatusCode)
	}
	LoginRequest, err := http.NewRequest("POST", hs.URL+"/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	rep, err = c.Do(LoginRequest)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		t.Fatal(err)
	}
	token = string(body)
	LoginStatusRequest, err = http.NewRequest("POST", hs.URL+"/loginstatus", nil)
	LoginStatusRequest.Header.Set(testHeaderName, token)
	if err != nil {
		t.Fatal(err)
	}
	rep, err = c.Do(LoginStatusRequest)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != http.StatusOK {
		t.Errorf("Status code error %d", rep.StatusCode)
	}
	TestRequest, err := http.NewRequest("POST", hs.URL+"/test", nil)
	TestRequest.Header.Set(testHeaderName, token)
	rep, err = c.Do(TestRequest)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != http.StatusOK {
		t.Errorf("HeaderMiddle status error %d", rep.StatusCode)
	}
	LogoutRequest, err := http.NewRequest("POST", hs.URL+"/logout", nil)
	LogoutRequest.Header.Set(testHeaderName, token)
	if err != nil {
		t.Fatal(err)
	}
	rep, err = c.Do(LogoutRequest)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	LoginStatusRequest, err = http.NewRequest("POST", hs.URL+"/loginstatus", nil)
	LoginStatusRequest.Header.Set(testHeaderName, token)
	if err != nil {
		t.Fatal(err)
	}
	rep, err = c.Do(LoginStatusRequest)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != http.StatusUnauthorized {
		t.Errorf("Status code error %d", rep.StatusCode)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	c = &http.Client{
		Jar: jar,
	}

	LoginStatusRequest, err = http.NewRequest("POST", hs.URL+"/cookie/loginstatus", nil)
	if err != nil {
		t.Fatal(err)
	}
	rep, err = c.Do(LoginStatusRequest)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != http.StatusUnauthorized {
		t.Errorf("Status code error %d", rep.StatusCode)
	}
	LoginRequest, err = http.NewRequest("POST", hs.URL+"/cookie/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	rep, err = c.Do(LoginRequest)
	if err != nil {
		t.Fatal(err)
	}
	body, err = ioutil.ReadAll(rep.Body)
	if err != nil {
		t.Fatal(err)
	}
	token = string(body)
	LoginStatusRequest, err = http.NewRequest("POST", hs.URL+"/cookie/loginstatus", nil)
	if err != nil {
		t.Fatal(err)
	}
	rep, err = c.Do(LoginStatusRequest)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != http.StatusOK {
		t.Errorf("Status code error %d", rep.StatusCode)
	}
	TestRequest, err = http.NewRequest("POST", hs.URL+"/cookie/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rep, err = c.Do(TestRequest)
	if err != nil {
		t.Fatal(err)
	}
	LogoutRequest, err = http.NewRequest("POST", hs.URL+"/cookie/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	rep, err = c.Do(LogoutRequest)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	LoginStatusRequest, err = http.NewRequest("POST", hs.URL+"/cookie/loginstatus", nil)
	if err != nil {
		t.Fatal(err)
	}
	rep, err = c.Do(LoginStatusRequest)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != http.StatusUnauthorized {
		t.Errorf("Status code error %d", rep.StatusCode)
	}
}

func TestCacheDriverPrefix(t *testing.T) {
	c := NewCacheDriver()
	p := "test"
	c.PrefixMode = PrefixModeEmpty
	t1, err := c.ConvertPrefix(p)
	if err != nil {
		panic(err)
	}
	if t1 != "" {
		t.Fatal(t1)
	}
	c.PrefixMode = PrefixModeRaw
	t1, err = c.ConvertPrefix(p)
	if err != nil {
		panic(err)
	}
	if t1 != p {
		t.Fatal(t1)
	}
	c.PrefixMode = PrefixModeMd5
	t1, err = c.ConvertPrefix(p)
	if err != nil {
		panic(err)
	}
	if len(t1) != 32 {
		t.Fatal(t1)
	}
	c.PrefixMode = PrefixModeEmpty
	c.Length = 10
	r, err := defaultTokenGenerater(c, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 10 {
		t.Fatal(r)
	}
	c.PrefixMode = ""
	c.Length = 10
	r, err = defaultTokenGenerater(c, p)
	t1, err = c.ConvertPrefix(p)
	if err != nil {
		panic(err)
	}
	if t1 != "" {
		t.Fatal(t1)
	}
}
