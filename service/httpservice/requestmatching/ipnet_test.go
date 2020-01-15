package requestmatching

import (
	"net/http"
	"testing"
)

func setRequestIP(r *http.Request, ip string) {
	r.RemoteAddr = ip + ":9999"
}
func setRequestIPv6(r *http.Request, ip string) {
	r.RemoteAddr = "[" + ip + "]:9999"
}
func TestIPNet(t *testing.T) {
	var r *http.Request
	r, _ = http.NewRequest("", "http://127.0.0.1/", nil)
	setRequestIP(r, "192.168.0.1")
	n := NewIPNets()
	if !mustMatch(n, r) {
		t.Fatal(n)
	}
	n.Add("127.0.0.1/24")
	if mustMatch(n, r) {
		t.Fatal(n)
	}
	n.Add("192.168.0.1/24")
	if !mustMatch(n, r) {
		t.Fatal(n)
	}
	r, _ = http.NewRequest("", "http://127.0.0.1/", nil)
	setRequestIPv6(r, "2001:0db8:0000:0000:0000:ff00:0042:8329")
	n = NewIPNets()
	if !mustMatch(n, r) {
		t.Fatal(n)
	}
	n.Add("2002::/64")
	if mustMatch(n, r) {
		t.Fatal(n)
	}
	n.Add("2001:0db8:0000:0000::/64")
	if !mustMatch(n, r) {
		t.Fatal(n)
	}

}

func TestContext(t *testing.T) {
	var r *http.Request
	r, _ = http.NewRequest("POST", "http://127.0.0.1/", nil)
	setRequestIP(r, "127.0.0.2")
	ipaddr := GetRequestIPAddress(r)
	if ipaddr != "127.0.0.2" {
		t.Fatal(ipaddr)
	}
	ip := GetRequestIP(r)
	if ip.String() != "127.0.0.2" {
		t.Fatal(ipaddr)
	}
	setRequestIP(r, "127.0.0.3")
	ipaddr = GetRequestIPAddress(r)
	if ipaddr != "127.0.0.2" {
		t.Fatal(ipaddr)
	}
	ip = GetRequestIP(r)
	if ip.String() != "127.0.0.2" {
		t.Fatal(ipaddr)
	}
}

func TestNilIP(t *testing.T) {
	var r *http.Request
	r, _ = http.NewRequest("", "http://127.0.0.1/", nil)
	setRequestIP(r, "192.168.0.1")
	n := NewIPNets()
	n.Add("192.168.0.1/32")
	if !mustMatch(n, r) {
		t.Fatal(n)
	}
	r, _ = http.NewRequest("", "http://127.0.0.1/", nil)
	r.RemoteAddr = "abcde"
	if mustMatch(n, r) {
		t.Fatal(n)
	}
}
