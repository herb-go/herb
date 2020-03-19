package session

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/herb-go/herb/cache"

	_ "github.com/herb-go/herb/cache/drivers/syncmapcache"
)

func getClientDriver(ttl time.Duration) *Store {
	s := MustClientStore([]byte("getClientDriver"), ttl)
	m, err := cache.NewMarshaler("json")
	if err != nil {
		panic(err)
	}
	s.Marshaler = m
	return s
}

func getTimeoutClientDriver(ttl time.Duration, UpdateActiveInterval time.Duration) *Store {
	s := MustClientStore([]byte("getTimeoutClientDriver"), ttl)
	m, err := cache.NewMarshaler("json")
	if err != nil {
		panic(err)
	}
	s.Marshaler = m
	s.UpdateActiveInterval = UpdateActiveInterval
	return s
}
func getBase64ClientDriver(ttl time.Duration) *Store {
	d := NewClientDriver()
	coc := NewClientDriverOptionConfig()
	coc.Key = []byte("getClientDriver")
	err := d.Init(coc)
	if err != nil {
		panic(err)
	}
	d.TokenMarshaler = func(s *ClientDriver, ts *Session) (err error) {
		var data []byte
		data, err = ts.Marshal()
		if err != nil {
			return err
		}
		ts.token = base64.StdEncoding.EncodeToString(data)
		return err
	}
	d.TokenUnmarshaler = func(s *ClientDriver, v *Session) (err error) {
		var data []byte
		data, err = base64.StdEncoding.DecodeString(v.token)
		if err != nil {
			return ErrDataNotFound
		}
		err = v.Unmarshal(v.token, data)
		if err != nil {
			return ErrDataNotFound
		}
		return nil

	}
	s := New()
	soc := NewOptionConfig()
	soc.Driver = d
	soc.TokenLifetime = ttl
	m, err := cache.NewMarshaler("json")
	if err != nil {
		panic(err)
	}
	s.Marshaler = m
	err = s.Init(soc)
	if err != nil {
		panic(err)
	}
	return s
}
func TestClientKey(t *testing.T) {
	s := getClientDriver(time.Hour)
	s2 := getTimeoutClientDriver(time.Hour, -1)
	model := "123456"
	var result string
	testKey := "testkey"
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
	td2 := s2.GetSession(td.MustToken())
	result = ""
	err = td2.Get(testKey, &result)
	if err != ErrDataNotFound {
		t.Fatal(err)
	}
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
}
func TestClientTD(t *testing.T) {
	var err error
	s := getClientDriver(time.Hour)
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

func TestClientRequest(t *testing.T) {
	var err error
	s := getClientDriver(time.Hour)
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
	if rep.Header.Get("set-cookie") == "" {
		t.Error("coolie update fail")
	}
	LogoutRequest, err := http.NewRequest("POST", hs.URL+"/cookie/logout", nil)
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

func TestClientTimeout(t *testing.T) {
	sforever := getTimeoutClientDriver(time.Hour, -1)
	s3second := getTimeoutClientDriver(3*time.Second, -1)
	s3secondwithAutoRefresh := getTimeoutClientDriver(3*time.Second, 1*time.Second)
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
