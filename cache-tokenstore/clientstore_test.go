package tokenstore

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	_ "github.com/herb-go/herb/cache/drivers/freecache"
)

func getClientStore(ttl time.Duration) Store {
	s := NewClientStore([]byte("getClientStore"), ttl)
	return s
}

func getTimeoutClientStore(ttl time.Duration, UpdateActiveInterval time.Duration) Store {
	s := NewClientStore([]byte("getTimeoutClientStore"), ttl)
	s.UpdateActiveInterval = UpdateActiveInterval
	return s
}

func TestClientStoreFieldInRequest(t *testing.T) {
	var err error
	s := getClientStore(-1)
	defer s.Close()
	model := "123456"
	modelAfterSet := "set"
	var result string
	testKey := "testkey"
	testOwner := "testowner"
	testHeaderName := "token"
	var token string
	field, err := s.RegisterField(testKey, &model)
	if err != nil {
		t.Fatal(err)
	}
	if field.Store != s {
		t.Errorf("Field store error")
	}
	if field.Type != reflect.TypeOf(model) {
		t.Errorf("Field type error")
	}
	var mux = http.NewServeMux()
	actionTest := func(w http.ResponseWriter, r *http.Request) {
		field.Get(r, &result)
		if result != model {
			t.Errorf("Field get error %s", result)
		}
		td, err := field.Store.GetRequestTokenData(r)
		if err != nil {
			t.Fatal(err)
		}
		tk, err := field.GetToken(r)
		if err != nil {
			t.Fatal(err)
		}
		if tk != token {
			t.Errorf("Field GetToken error %s", tk)
		}
		result = ""
		err = field.GetFromToken(td.MustToken(), &result)
		if err != nil {
			t.Fatal(err)
		}
		if result != model {
			t.Errorf("Field GetFromToken error %s", tk)
		}
		ex, err := field.ExpiredAt(r)
		if err != nil {
			t.Fatal(err)
		}
		if ex > 0 {
			t.Errorf("Field ExpiredAt error %d", ex)
		}
		mutex, err := field.RwMutex(r)
		if err != nil {
			t.Fatal(err)
		}
		if mutex != td.Mutex {
			t.Errorf("Field mutex error")
		}
		err = field.Set(r, modelAfterSet)
		if err != nil {
			t.Fatal(err)
		}
		result = ""
		err = field.Get(r, &result)
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
	actionLogin := func(w http.ResponseWriter, r *http.Request) {
		td := field.MustLogin(r, testOwner, model)
		w.Write([]byte(td.MustToken()))
	}
	actionHeaderLogin := func(w http.ResponseWriter, r *http.Request) {
		s.HeaderMiddleware(testHeaderName)(w, r, actionLogin)
	}
	mux.HandleFunc("/login", actionHeaderLogin)
	mux.HandleFunc("/test", actionHeaderTest)
	hs := httptest.NewServer(mux)
	c := &http.Client{}
	LoginRequest, err := http.NewRequest("POST", hs.URL+"/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	rep, err := c.Do(LoginRequest)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		t.Fatal(err)
	}
	token = string(body)
	TestRequest, err := http.NewRequest("POST", hs.URL+"/test", nil)
	TestRequest.Header.Set(testHeaderName, token)
	rep, err = c.Do(TestRequest)
	if err != nil {
		t.Fatal(err)
	}
	if rep.StatusCode != 200 {
		t.Errorf("HeaderMiddle status error %d", rep.StatusCode)
	}
}

func TestClientStoreTimeout(t *testing.T) {
	sforever := getTimeoutClientStore(-1, -1)
	s3second := getTimeoutClientStore(3*time.Second, -1)
	s3secondwithAutoRefresh := getTimeoutClientStore(3*time.Second, 1*time.Second)
	testOwner := "testowner"
	model := "123456"
	var result string
	testKey := "testkey"
	fieldForever, err := sforever.RegisterField(testKey, &model)
	if err != nil {
		panic(err)
	}
	tdForeverKey, err := sforever.GenerateToken(testOwner)
	if err != nil {
		panic(err)
	}
	tdForever, err := sforever.GenerateTokenData(tdForeverKey)
	if err != nil {
		panic(err)
	}
	err = fieldForever.SaveTo(tdForever, model)
	if err != nil {
		panic(err)
	}
	field3second, err := s3second.RegisterField(testKey, &model)
	if err != nil {
		panic(err)
	}
	td3secondKey, err := s3second.GenerateToken(testOwner)
	if err != nil {
		panic(err)
	}
	td3second, err := s3second.GenerateTokenData(td3secondKey)
	if err != nil {
		panic(err)
	}
	err = field3second.SaveTo(td3second, model)
	if err != nil {
		panic(err)
	}
	field3secondwithAutoRefresh, err := s3secondwithAutoRefresh.RegisterField(testKey, &model)
	if err != nil {
		panic(err)
	}
	td3secondwithAutoRefreshKey, err := s3secondwithAutoRefresh.GenerateToken(testOwner)
	if err != nil {
		panic(err)
	}
	td3secondwithAutoRefresh, err := s3secondwithAutoRefresh.GenerateTokenData(td3secondwithAutoRefreshKey)
	if err != nil {
		panic(err)
	}
	err = field3secondwithAutoRefresh.SaveTo(td3secondwithAutoRefresh, model)
	if err != nil {
		panic(err)
	}
	tdForever.Save()
	td3second.Save()
	td3secondwithAutoRefresh.Save()
	time.Sleep(2 * time.Second)
	tdForever, err = sforever.GetTokenData(tdForever.MustToken())
	if err != nil {
		panic(err)
	}
	td3second, err = s3second.GetTokenData(td3second.MustToken())
	if err != nil {
		panic(err)
	}
	td3secondwithAutoRefresh, err = s3secondwithAutoRefresh.GetTokenData(td3secondwithAutoRefresh.MustToken())
	if err != nil {
		panic(err)
	}
	result = ""
	err = fieldForever.LoadFrom(tdForever, &result)
	if result != model {
		t.Errorf("Timeout error %s", result)
	}
	result = ""
	err = field3second.LoadFrom(td3second, &result)
	if result != model {
		t.Errorf("Timeout error %s", result)
	}
	result = ""
	err = field3secondwithAutoRefresh.LoadFrom(td3secondwithAutoRefresh, &result)
	if result != model {
		t.Errorf("Timeout error %s", result)
	}
	tdForever.Save()
	td3second.Save()
	td3secondwithAutoRefresh.Save()
	time.Sleep(2 * time.Second)
	tdForever, err = sforever.GetTokenData(tdForever.MustToken())
	if err != nil {
		panic(err)
	}
	td3second, err = s3second.GetTokenData(td3second.MustToken())
	if err != ErrDataNotFound {
		panic(err)
	}
	td3secondwithAutoRefresh, err = s3secondwithAutoRefresh.GetTokenData(td3secondwithAutoRefresh.MustToken())
	if err != nil {
		panic(err)
	}
	result = ""
	err = fieldForever.LoadFrom(tdForever, &result)
	if result != model {
		t.Errorf("Timeout error %s", result)
	}
	result = ""
	err = field3second.LoadFrom(td3second, &result)
	if err != ErrDataNotFound {
		t.Errorf("Timeout error %s", err)
	}
	result = ""
	err = field3secondwithAutoRefresh.LoadFrom(td3secondwithAutoRefresh, &result)
	if result != model {
		t.Errorf("Timeout error %s", result)
	}
	tdForever.Save()
	td3second.Save()
	td3secondwithAutoRefresh.Save()
	time.Sleep(4 * time.Second)
	tdForever, err = sforever.GetTokenData(tdForever.MustToken())
	if err != nil {
		panic(err)
	}
	td3second, err = s3second.GetTokenData(td3second.MustToken())
	if err != ErrDataNotFound {
		panic(err)
	}
	td3secondwithAutoRefresh, err = s3secondwithAutoRefresh.GetTokenData(td3secondwithAutoRefresh.MustToken())
	if err != ErrDataNotFound {
		panic(err)
	}
	result = ""
	err = fieldForever.LoadFrom(tdForever, &result)
	if result != model {
		t.Errorf("Timeout error %s", result)
	}
	result = ""
	err = field3second.LoadFrom(td3second, &result)
	if err != ErrDataNotFound {
		t.Errorf("Timeout error %s", err)
	}
	result = ""
	err = field3secondwithAutoRefresh.LoadFrom(td3secondwithAutoRefresh, &result)
	if err != ErrDataNotFound {
		t.Errorf("Timeout error %s", err)
	}
}
