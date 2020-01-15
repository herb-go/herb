package requestmatching

import (
	"net/http"
	"testing"
)

func mustMatch(p Pattern, r *http.Request) bool {
	result, err := p.MatchRequest(r)
	if err != nil {
		panic(err)
	}
	return result
}

func TestPatternConfig(t *testing.T) {
	var c PatternConfig
	var p Pattern
	var err error
	p, err = c.CreatePattern()
	if err != nil {
		t.Fatal(err)
	}
	if !mustMatch(p, nil) {
		t.Fatal()
	}
	c.IPList = []string{"127.0.0.1"}
	p, err = c.CreatePattern()
	if p != nil || err == nil {
		t.Fatal(p, err)
	}
	c = PatternConfig{
		IPList:     []string{"127.0.0.1/8", "127.0.0.2/24"},
		URLList:    []string{"127.0.0.1/path", "path2", "path3"},
		PrefixList: []string{"127.0.0.1/path", "path2", "path3", "path4"},
		MethodList: []string{"get", "post", "3", "4", "5"},
		ExtList:    []string{".html", "", ".3", ".4", ".5", ".6"},
		Disabled:   true,
		Not:        true,
		And:        true,
		Patterns:   []*PatternConfig{&PatternConfig{}},
	}
	p, err = c.CreatePattern()
	if p == nil || err != nil {
		t.Fatal(p, err)
	}
	pattern := p.(*PlainPattern)
	if pattern == nil {
		t.Fatal(p)
	}
	if len(*pattern.IPNets) != 2 ||
		len(*pattern.Paths) != 3 ||
		len(*pattern.Prefixs) != 4 ||
		len(*pattern.Methods) != 5 ||
		len(*pattern.Exts) != 6 ||
		pattern.Disabled != true ||
		pattern.Not != true ||
		pattern.And != true ||
		len(pattern.Patterns) != 1 {
		t.Fatal(pattern)
	}
}

func TestPattern(t *testing.T) {
	var p = NewPlainPattern()
	var r *http.Request
	if mustMatch(p, nil) != true {
		t.Fatal(p)
	}
	p.Disabled = true
	if mustMatch(p, nil) != false {
		t.Fatal(p)
	}
	p = NewPlainPattern()
	r, _ = http.NewRequest("POST", "http://127.0.0.1/abc", nil)
	setRequestIP(r, "127.0.0.1")
	p.IPNets.Add("127.0.0.2/32")
	if mustMatch(p, r) != false {
		t.Fatal(p)
	}
	p.IPNets.Add("127.0.0.1/32")
	if mustMatch(p, r) != true {
		t.Fatal(p)
	}
	p = NewPlainPattern()
	r, _ = http.NewRequest("POST", "http://127.0.0.1/abc", nil)
	p.Methods.Add("get")
	if mustMatch(p, r) != false {
		t.Fatal(p)
	}
	p.Methods.Add("Post")
	if mustMatch(p, r) != true {
		t.Fatal(p)
	}
	p = NewPlainPattern()
	r, _ = http.NewRequest("POST", "http://127.0.0.1/abc", nil)
	p.Exts.Add(".html")
	if mustMatch(p, r) != false {
		t.Fatal(p)
	}
	p.Exts.Add("")
	if mustMatch(p, r) != true {
		t.Fatal(p)
	}
	p = NewPlainPattern()
	r, _ = http.NewRequest("POST", "http://127.0.0.1/abc?dd=15", nil)
	p.Paths.Add("127.0.0.2/abc")
	if mustMatch(p, r) != false {
		t.Fatal(p)
	}
	p.Paths.Add("/abc")
	if mustMatch(p, r) != true {
		t.Fatal(p)
	}
	p = NewPlainPattern()
	r, _ = http.NewRequest("POST", "http://127.0.0.1/abc", nil)
	p.Prefixs.Add("127.0.0.2/abc")
	if mustMatch(p, r) != false {
		t.Fatal(p)
	}
	p.Prefixs.Add("/a")
	if mustMatch(p, r) != true {
		t.Fatal(p)
	}

	p = NewPlainPattern()
	r, _ = http.NewRequest("POST", "http://127.0.0.1/abc", nil)
	setRequestIP(r, "127.0.0.1")
	p.IPNets.Add("127.0.0.2/32")
	if mustMatch(p, r) != false {
		t.Fatal(p)
	}
	p.IPNets.Add("127.0.0.1/32")
	if mustMatch(p, r) != true {
		t.Fatal(p)
	}
	p.Paths.Add("/cba")
	if mustMatch(p, r) != false {
		t.Fatal(p)
	}
	p.Paths.Add("/abc")
	if mustMatch(p, r) != true {
		t.Fatal(p)
	}
}

func createSuccessPattern() *PlainPattern {
	s := NewPlainPattern()
	s.Paths.Add("/success")

	return s
}

func createFailPattern() *PlainPattern {
	f := NewPlainPattern()
	f.Paths.Add("/fail")
	return f
}

func createSuccessRequest() *http.Request {
	r, _ := http.NewRequest("", "/success", nil)
	return r
}
func TestData(t *testing.T) {
	s := createSuccessPattern()
	f := createFailPattern()
	r := createSuccessRequest()
	if !mustMatch(s, r) || mustMatch(f, r) {
		t.Fatal(s, f, r)
	}
}

func TestSubPattern(t *testing.T) {
	var p *PlainPattern
	r := createSuccessRequest()

	p = createSuccessPattern()
	p.And = true
	p.Patterns = append(p.Patterns, createFailPattern())
	if mustMatch(p, r) {
		t.Fatal()
	}
	p = createSuccessPattern()
	p.And = true
	p.Patterns = append(p.Patterns, createSuccessPattern())
	if !mustMatch(p, r) {
		t.Fatal()
	}

	p = createSuccessPattern()
	p.And = false
	p.Patterns = append(p.Patterns, createFailPattern())
	if !mustMatch(p, r) {
		t.Fatal()
	}
	p = createSuccessPattern()
	p.And = true
	p.Patterns = append(p.Patterns, createSuccessPattern())
	if !mustMatch(p, r) {
		t.Fatal()
	}
	p = createFailPattern()
	p.And = true
	p.Patterns = append(p.Patterns, createSuccessPattern())
	if !mustMatch(p, r) {
		t.Fatal()
	}
}

func TestOperationPattern(t *testing.T) {
	var result bool
	var err error
	r := createSuccessRequest()

	p := createSuccessPattern()
	p.Not = true
	if mustMatch(p, r) {
		t.Fatal()
	}
	result, err = MatchAll(r, createSuccessPattern(), createSuccessPattern())
	if !result || err != nil {
		t.Fatal(result, err)
	}
	result, err = MatchAll(r, createSuccessPattern(), createFailPattern())
	if result || err != nil {
		t.Fatal(result, err)
	}
	result, err = MatchAny(r, createSuccessPattern(), createSuccessPattern())
	if !result || err != nil {
		t.Fatal(result, err)
	}
	result, err = MatchAny(r, createSuccessPattern(), createFailPattern())
	if !result || err != nil {
		t.Fatal(result, err)
	}
	result, err = MatchAny(r, createFailPattern(), createFailPattern())
	if result || err != nil {
		t.Fatal(result, err)
	}
}
