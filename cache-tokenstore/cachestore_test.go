package tokenstore

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/herb-go/herb/cache"
	_ "github.com/herb-go/herb/cache/drivers/freecache"
)

func getStore(ttl time.Duration) Store {
	c := cache.New()
	err := c.OpenJSON([]byte(testCache))
	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	s := New(c, ttl)
	return s
}

func getTimeoutStore(ttl time.Duration, UpdateActiveInterval time.Duration) Store {
	c := cache.New()
	err := c.OpenJSON([]byte(testCache))
	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	s := New(c, ttl)
	s.UpdateActiveInterval = UpdateActiveInterval
	return s
}

func TestField(t *testing.T) {
	var err error
	s := getStore(-1)
	defer s.Close()
	model := "123456"
	var result string
	testKey := "testkey"
	type modelStruct struct {
		data string
	}
	structModel := modelStruct{
		data: "test",
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
	_, err = s.RegisterField(testKey, model)
	if err != ErrMustRegistePtr {
		t.Fatal(err)
	}
	_, err = s.RegisterField(testKey, nil)
	if err != ErrNilPointer {
		t.Fatal(err)
	}
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
	td := field.MustLoginTokenData(testOwner, model)
	result = ""
	err = field.LoadFrom(td, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result != model {
		t.Fatal(err)
	}
	result = ""
	err = field.GetFromToken(td.MustToken(), &result)
	if err != nil {
		t.Fatal(err)
	}
	if result != model {
		t.Errorf("Field GetFromToken error")
	}
	fieldStruct, err := s.RegisterField(testIntKey, &modelStruct{})
	if err != nil {
		t.Fatal(err)
	}
	err = fieldStruct.SaveTo(td, structModel)
	if err != nil {
		t.Fatal(err)
	}
	err = fieldStruct.LoadFrom(td, &resutStruct)
	if err != nil {
		t.Fatal(err)
	}
	if resutStruct.data != structModel.data {
		t.Errorf("field Struct error %s", resutStruct.data)
	}

	fieldInt, err := s.RegisterField(testStructKey, &resultInt)
	if err != nil {
		t.Fatal(err)
	}
	err = fieldInt.SaveTo(td, modelInt)
	if err != nil {
		t.Fatal(err)
	}
	err = fieldInt.LoadFrom(td, &resultInt)
	if err != nil {
		t.Fatal(err)
	}
	if resultInt != modelInt {
		t.Errorf("field int error %d", resultInt)
	}
	fieldBytes, err := s.RegisterField(testBytesKey, &resultBytes)
	if err != nil {
		t.Fatal(err)
	}
	err = fieldBytes.SaveTo(td, modelBytes)
	if err != nil {
		t.Fatal(err)
	}
	err = fieldBytes.LoadFrom(td, &resultBytes)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(resultBytes, modelBytes) != 0 {
		t.Errorf("field Bytes error %s", string(resultBytes))
	}

	fieldMap, err := s.RegisterField(testMapKey, &resultMap)
	if err != nil {
		t.Fatal(err)
	}
	err = fieldMap.SaveTo(td, modelMap)
	if err != nil {
		t.Fatal(err)
	}
	err = fieldMap.LoadFrom(td, &resultMap)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(resultMap, modelMap) {
		t.Error("field Maps error")
	}
}

func TestFieldInRequest(t *testing.T) {
	var err error
	s := getStore(-1)
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
		s.HeaderMiddleware(testHeaderName)(w, r, func(w http.ResponseWriter, r *http.Request) {
			field.Get(r, &result)
			if result != model {
				t.Errorf("Field get error %s", result)
			}
			td, err := field.Store.GetRequestTokenData(r)
			if err != nil {
				t.Fatal(err)
			}
			if token != td.MustToken() {
				t.Errorf("field.Store.GetRequestTokenData error %s", token)
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
		})
	}
	actionLogin := func(w http.ResponseWriter, r *http.Request) {
		s.HeaderMiddleware(testHeaderName)(w, r, func(w http.ResponseWriter, r *http.Request) {
			td := field.MustLogin(r, testOwner, model)
			w.Write([]byte(td.MustToken()))
		})
	}
	mux.HandleFunc("/login", actionLogin)
	mux.HandleFunc("/test", actionTest)
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

func TestTokenDataMarshal(t *testing.T) {
	var err error
	testOwner := "testowner"
	model := "123456"
	var result string
	testKey := "testkey"
	testKey2 := "testkey2"
	testToken := "testtoken"
	s := getStore(-1)
	defer s.Close()
	field, err := s.RegisterField(testKey, &model)
	if err != nil {
		panic(err)
	}
	td, err := s.GenerateTokenData(testOwner)
	if err != nil {
		panic(err)
	}
	err = field.SaveTo(td, model)
	if err != nil {
		panic(err)
	}
	bytes, err := td.Marshal()
	if err != nil {
		panic(err)
	}
	td2 := NewTokenData(testToken, s)
	err = td2.Unmarshal(testKey2, bytes)
	if err != nil {
		panic(err)
	}
	err = field.LoadFrom(td2, &result)
	if err != nil {
		panic(err)
	}
	if result != model {
		t.Errorf("Tokendata Unmarshal err %s", result)
	}

}
func TestTokenData(t *testing.T) {
	var err error
	s := getStore(-1)
	defer s.Close()
	testOwner := "testowner"
	model := "123456"
	var result string
	testKey := "testkey"
	field, err := s.RegisterField(testKey, &model)
	if err != nil {
		panic(err)
	}
	td, err := s.GenerateTokenData(testOwner)
	if err != nil {
		panic(err)
	}
	err = field.SaveTo(td, model)
	if err != nil {
		panic(err)
	}
	testToken := td.MustToken()
	td2 := NewTokenData(testToken, s)
	err = td2.Load()
	if err != ErrDataNotFound {
		t.Fatal(err)
	}
	err = field.LoadFrom(td2, &result)
	if err != ErrDataNotFound {
		t.Fatal(err)
	}
	err = td.Save()
	if err != nil {
		panic(err)
	}

	td3 := NewTokenData(testToken, s)
	err = td3.Load()
	if err != nil {
		panic(err)
	}
	result = ""
	err = field.LoadFrom(td3, &result)
	if err != nil {
		panic(err)
	}
	if result != model {
		t.Errorf("Tokendata save/load error %s", result)
	}
	err = td.DeleteAndSave()
	if err != nil {
		panic(err)
	}
	td4 := NewTokenData(testToken, s)
	err = td4.Load()
	if err != ErrDataNotFound {
		t.Fatal(err)
	}
	err = field.LoadFrom(td4, &result)
	if err != ErrDataNotFound {
		t.Fatal(err)
	}
}

func TestTimeout(t *testing.T) {
	sforever := getTimeoutStore(-1, -1)
	s3second := getTimeoutStore(3*time.Second, -1)
	s3secondwithAutoRefresh := getTimeoutStore(3*time.Second, 1*time.Second)
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
	if err != ErrDataNotFound {
		t.Errorf("Timeout error %s", err)
	}
	result = ""
	err = field3secondwithAutoRefresh.LoadFrom(td3secondwithAutoRefresh, &result)
	if err != ErrDataNotFound {
		t.Errorf("Timeout error %s", err)
	}
}