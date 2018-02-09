package forwarded

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/herb-go/herb/middleware"
)

const xForwardedForHeader = "X-Forwarded-For"
const xForwardedHostHeader = "X-Forwarded-Host"
const xForwardedProtoHeader = "X-Forwarded-Proto"

var s *httptest.Server

func doRequest(t *testing.T, header http.Header) (map[string]string, int) {
	var url = "http://" + s.Listener.Addr().String() + "/"
	var req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range header {
		for _, value := range v {
			req.Header.Add(k, value)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var data = map[string]string{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Fatal(err)
	}
	return data, resp.StatusCode
}
func TestMiddleware(t *testing.T) {
	var app = middleware.New()
	var data map[string]string
	var statusCode int
	var middleware = Middleware{}
	app.
		Use(middleware.ServeMiddleware).
		HandleFunc(func(w http.ResponseWriter, r *http.Request) {
			var result = map[string]string{}
			result["RemoteAddr"] = r.RemoteAddr
			result["Host"] = r.Host
			result["Scheme"] = r.URL.Scheme
			var data, err = json.Marshal(result)
			if err != nil {
				panic(err)
			}
			_, err = w.Write(data)
			if err != nil {
				panic(err)
			}
		})
	s = httptest.NewServer(app)
	defer s.Close()
	var host = s.Listener.Addr().String()
	var addr = "127.0.0.1"
	var scheme = ""
	data, statusCode = doRequest(t, nil)
	if statusCode != http.StatusOK {
		t.Fatal(statusCode)
	}
	if strings.Split(data["RemoteAddr"], ":")[0] != addr {
		t.Error(data["RemoteAddr"])
	}
	if data["Host"] != host {
		t.Error(data["Host"])
	}
	if data["scheme"] != scheme {
		t.Error(data["scheme"])
	}

	middleware.Status = 1
	data, statusCode = doRequest(t, nil)

	if statusCode != http.StatusOK {
		t.Fatal(statusCode)
	}
	data, statusCode = doRequest(t, nil)
	if strings.Split(data["RemoteAddr"], ":")[0] != addr {
		t.Error(data["RemoteAddr"])
	}
	if data["Host"] != host {
		t.Error(data["Host"])
	}
	if data["scheme"] != scheme {
		t.Error(data["scheme"])
	}

	middleware.ForwardedForHeader = xForwardedForHeader
	data, statusCode = doRequest(t, nil)
	if statusCode != http.StatusOK {
		t.Fatal(statusCode)
	}
	if strings.Split(data["RemoteAddr"], ":")[0] != addr {
		t.Error(data["RemoteAddr"])
	}

	var forwardedIP = "1234567"
	data, statusCode = doRequest(t, map[string][]string{
		xForwardedForHeader: []string{forwardedIP},
	})
	if statusCode != http.StatusOK {
		t.Fatal(statusCode)
	}
	if strings.Split(data["RemoteAddr"], ":")[0] != forwardedIP {
		t.Error(data["RemoteAddr"])
	}

	middleware.ForwardedHostHeader = xForwardedHostHeader
	data, statusCode = doRequest(t, nil)
	if statusCode != http.StatusOK {
		t.Fatal(statusCode)
	}
	if data["Host"] != host {
		t.Error(data["Host"])
	}

	var forwardedHost = "github.com"
	data, statusCode = doRequest(t, map[string][]string{
		xForwardedHostHeader: []string{forwardedHost},
	})
	if statusCode != http.StatusOK {
		t.Fatal(statusCode)
	}
	if data["Host"] != forwardedHost {
		t.Error(data["Host"])
	}

	middleware.ForwardedProtoHeader = xForwardedProtoHeader
	data, statusCode = doRequest(t, nil)
	if statusCode != http.StatusOK {
		t.Fatal(statusCode)

	}
	if data["Scheme"] != "" {
		t.Error(data["Scheme"])
	}
	var forwardedScheme = "https"
	data, statusCode = doRequest(t, map[string][]string{
		xForwardedProtoHeader: []string{forwardedScheme},
	})
	if statusCode != http.StatusOK {
		t.Fatal(statusCode)
	}
	if data["Scheme"] != forwardedScheme {
		t.Error(data["Scheme"])
	}

	var forwardedToken = "X-Forwarded-Token"
	var forwardedTokenWrong = "wrongtoken"
	var forwardedValue = "value"
	var forwardedValueWrong = "wrongvalue"

	middleware.ForwardedTokenHeader = forwardedToken
	data, statusCode = doRequest(t, nil)
	if statusCode != http.StatusBadRequest {
		t.Fatal(statusCode)
	}

	data, statusCode = doRequest(t, map[string][]string{
		forwardedTokenWrong: []string{forwardedValue},
	})
	middleware.ForwardedTokenValue = forwardedValue

	if statusCode != http.StatusBadRequest {
		t.Fatal(statusCode)
	}
	data, statusCode = doRequest(t, map[string][]string{
		forwardedToken: []string{forwardedValueWrong},
	})
	if statusCode != http.StatusBadRequest {
		t.Fatal(statusCode)
	}
	data, statusCode = doRequest(t, map[string][]string{
		forwardedToken: []string{forwardedValue},
	})
	if statusCode != http.StatusOK {
		t.Error(middleware)
		t.Fatal(statusCode)
	}

	middleware.SetTokenFailedAction(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
	data, statusCode = doRequest(t, nil)
	if statusCode != http.StatusForbidden {
		t.Fatal(statusCode)
	}
	middleware.Status = 0
	data, statusCode = doRequest(t, nil)

	if statusCode != http.StatusOK {
		t.Fatal(statusCode)
	}
	data, statusCode = doRequest(t, nil)
	if strings.Split(data["RemoteAddr"], ":")[0] != addr {
		t.Error(data["RemoteAddr"])
	}
	if data["Host"] != host {
		t.Error(data["Host"])
	}
	if data["scheme"] != scheme {
		t.Error(data["scheme"])
	}
}
