package fetch

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type testServerData struct {
	Name string
}

type testServerResp struct {
	testServerData
	Header string
}

func testJSONAction(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	data := testServerData{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}
	output, err := json.Marshal(
		testServerResp{
			data,
			r.Header.Get("test"),
		},
	)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(output)
	if err != nil {
		panic(err)
	}
}

func testXMLAction(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	data := testServerData{}
	err = xml.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}
	output, err := xml.Marshal(
		testServerResp{
			data,
			r.Header.Get("test"),
		},
	)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(output)
	if err != nil {
		panic(err)
	}
}

func testPostAction(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	output := strings.Join([]string{string(content), r.Header.Get("test")}, ",")
	_, err = w.Write([]byte(output))
	if err != nil {
		panic(err)
	}
}

func testGetAction(w http.ResponseWriter, r *http.Request) {
	output := strings.Join([]string{r.URL.Query().Get("Name"), r.Header.Get("test")}, ",")
	_, err := w.Write([]byte(output))
	if err != nil {
		panic(err)
	}
}
func TestServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/json", testJSONAction)
	mux.HandleFunc("/xml", testXMLAction)
	mux.HandleFunc("/post", testPostAction)
	mux.HandleFunc("/get", testGetAction)

	ts := httptest.NewServer(mux)
	defer ts.Close()
	server := Server{
		Host: ts.URL,
		Headers: http.Header{
			"test": []string{
				"testheader",
			},
		},
	}
	APIJSON := server.EndPoint("POST", "/json")
	APIXML := server.EndPoint("POST", "/xml")
	APIPost := server.EndPoint("POST", "/post")
	APIGet := server.EndPoint("Get", "/get")

	req, err := APIJSON.NewJSONRequest(nil, testServerData{Name: "testname"})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	result := &testServerResp{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Header != "testheader" {
		t.Error(result.Header)
	}
	if result.Name != "testname" {
		t.Error(result.Name)
	}

	req, err = APIXML.NewXMLRequest(nil, testServerData{Name: "testname"})
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	result = &testServerResp{}
	err = xml.Unmarshal(content, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Header != "testheader" {
		t.Error(result.Header)
	}
	if result.Name != "testname" {
		t.Error(result.Name)
	}

	req, err = APIPost.NewRequest(nil, []byte("testname"))
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	resultbody := strings.Split(string(content), ",")
	if resultbody[1] != "testheader" {
		t.Error(result.Header)
	}
	if resultbody[0] != "testname" {
		t.Error(result.Name)
	}

	req, err = APIGet.NewRequest(url.Values{"Name": []string{"testname"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	resultbody = strings.Split(string(content), ",")
	if resultbody[1] != "testheader" {
		t.Error(result.Header)
	}
	if resultbody[0] != "testname" {
		t.Error(result.Name)
	}
}
