package requestmatching

import (
	"net/http"
	"testing"
)

func TestPatternConfig(t *testing.T) {
	var c PatternConfig
	var p Pattern
	var err error
	p, err = c.CreatePattern()
	if err != nil {
		t.Fatal(err)
	}
	if !MustMatch(nil, p) {
		t.Fatal()
	}
	c.IPList = []string{"127.0.0.1"}
	p, err = c.CreatePattern()
	if p != nil || err == nil {
		t.Fatal(p, err)
	}
	c = PatternConfig{
		IPList:          []string{"127.0.0.1/8", "127.0.0.2/24"},
		URLList:         []string{"127.0.0.1/path", "path2", "path3"},
		PrefixList:      []string{"127.0.0.1/path", "path2", "path3", "path4"},
		MethodList:      []string{"get", "post", "3", "4", "5"},
		ExtList:         []string{".html", "", ".3", ".4", ".5", ".6"},
		SuffixList:      []string{"127.0.0.1/path", "path2", "path3", "path4", "path5", "path6", "path7"},
		KeywordList:     []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		RegExpList:      []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
		HeaderList:      []string{"f:1", "f:2", "f:3", "f:4", "f:5", "f:6", "f:7", "f:8", "f:9", "f:10"},
		ContentTypeList: []string{"application/x-www-form-urlencoded", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
		Disabled:        true,
		Not:             true,
		And:             true,
		Patterns:        []*PatternConfig{&PatternConfig{}},
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
		len(pattern.Patterns) != 1 ||
		len(*pattern.Suffixs) != 7 ||
		len(*pattern.Keywords) != 8 ||
		len(*pattern.RegExps) != 9 ||
		len(*pattern.Headers) != 10 ||
		len(*pattern.ContentTypes) != 11 {
		t.Fatal(pattern)
	}
}

func TestPattern(t *testing.T) {
	var p = NewPlainPattern()
	var r *http.Request
	if MustMatch(nil, p) != true {
		t.Fatal(p)
	}
	p.Disabled = true
	if MustMatch(nil, p) != false {
		t.Fatal(p)
	}
	p = NewPlainPattern()
	r, _ = http.NewRequest("POST", "http://127.0.0.1/abc", nil)
	setRequestIP(r, "127.0.0.1")
	p.IPNets.Add("127.0.0.2/32")
	if MustMatch(r, p) != false {
		t.Fatal(p)
	}
	p.IPNets.Add("127.0.0.1/32")
	if MustMatch(r, p) != true {
		t.Fatal(p)
	}
	p = NewPlainPattern()
	r, _ = http.NewRequest("POST", "http://127.0.0.1/abc", nil)
	p.Methods.Add("get")
	if MustMatch(r, p) != false {
		t.Fatal(p)
	}
	p.Methods.Add("Post")
	if MustMatch(r, p) != true {
		t.Fatal(p)
	}
	p = NewPlainPattern()
	r, _ = http.NewRequest("POST", "http://127.0.0.1/abc", nil)
	p.Exts.Add(".html")
	if MustMatch(r, p) != false {
		t.Fatal(p)
	}
	p.Exts.Add("")
	if MustMatch(r, p) != true {
		t.Fatal(p)
	}
	p = NewPlainPattern()
	r, _ = http.NewRequest("POST", "http://127.0.0.1/abc?dd=15", nil)
	p.Paths.Add("127.0.0.2/abc")
	if MustMatch(r, p) != false {
		t.Fatal(p)
	}
	p.Paths.Add("/abc")
	if MustMatch(r, p) != true {
		t.Fatal(p)
	}
	p = NewPlainPattern()
	r, _ = http.NewRequest("POST", "http://127.0.0.1/abc", nil)
	p.Prefixs.Add("127.0.0.2/abc")
	if MustMatch(r, p) != false {
		t.Fatal(p)
	}
	p.Prefixs.Add("/a")
	if MustMatch(r, p) != true {
		t.Fatal(p)
	}

	p = NewPlainPattern()
	r, _ = http.NewRequest("POST", "http://127.0.0.1/abc", nil)
	setRequestIP(r, "127.0.0.1")
	p.IPNets.Add("127.0.0.2/32")
	if MustMatch(r, p) != false {
		t.Fatal(p)
	}
	p.IPNets.Add("127.0.0.1/32")
	if MustMatch(r, p) != true {
		t.Fatal(p)
	}
	p.Paths.Add("/cba")
	if MustMatch(r, p) != false {
		t.Fatal(p)
	}
	p.Paths.Add("/abc")
	if MustMatch(r, p) != true {
		t.Fatal(p)
	}
}

var patternConfigSuccess = &PatternConfig{
	URLList: []string{"/success"},
}
var patternConfigFail = &PatternConfig{
	URLList: []string{"/fail"},
}

func createSuccessPattern() *PlainPattern {

	return MustCreatePattern(patternConfigSuccess).(*PlainPattern)
}

func createFailPattern() *PlainPattern {
	return MustCreatePattern(patternConfigFail).(*PlainPattern)

}

func createSuccessRequest() *http.Request {
	r, _ := http.NewRequest("", "/success", nil)
	return r
}
func TestData(t *testing.T) {
	s := createSuccessPattern()
	f := createFailPattern()
	r := createSuccessRequest()
	if !MustMatch(r, s) || MustMatch(r, f) {
		t.Fatal(s, f, r)
	}
}

func TestSubPattern(t *testing.T) {
	var p *PlainPattern
	r := createSuccessRequest()

	p = createSuccessPattern()
	p.And = true
	p.Patterns = append(p.Patterns, createFailPattern())
	if MustMatch(r, p) {
		t.Fatal()
	}
	p = createSuccessPattern()
	p.And = true
	p.Patterns = append(p.Patterns, createSuccessPattern())
	if !MustMatch(r, p) {
		t.Fatal()
	}

	p = createSuccessPattern()
	p.And = false
	p.Patterns = append(p.Patterns, createFailPattern())
	if !MustMatch(r, p) {
		t.Fatal()
	}
	p = createSuccessPattern()
	p.And = true
	p.Patterns = append(p.Patterns, createSuccessPattern())
	if !MustMatch(r, p) {
		t.Fatal()
	}
	p = createFailPattern()
	p.And = true
	p.Patterns = append(p.Patterns, createSuccessPattern())
	if !MustMatch(r, p) {
		t.Fatal()
	}
}

func TestOperationPattern(t *testing.T) {
	var result bool
	var err error
	r := createSuccessRequest()

	p := createSuccessPattern()
	p.Not = true
	if MustMatch(r, p) {
		t.Fatal()
	}
	result = MustMatchAll(r)
	if !result {
		t.Fatal(result)
	}
	result, err = MatchAll(r, createSuccessPattern(), createSuccessPattern())
	if !result || err != nil {
		t.Fatal(result, err)
	}
	result = MustMatchAll(r, createSuccessPattern(), createFailPattern())
	if result {
		t.Fatal(result)
	}
	result = MustMatchAny(r)
	if !result {
		t.Fatal(result)
	}
	result, err = MatchAny(r, createSuccessPattern(), createSuccessPattern())
	if !result || err != nil {
		t.Fatal(result, err)
	}
	result, err = MatchAny(r, createSuccessPattern(), createFailPattern())
	if !result || err != nil {
		t.Fatal(result, err)
	}
	result = MustMatchAny(r, createFailPattern(), createFailPattern())
	if result {
		t.Fatal(result)
	}
}

func TestConfigs(t *testing.T) {
	configAllSuccess := &PatternAllConfig{patternConfigSuccess, patternConfigSuccess}
	configAllSuccessAndFail := &PatternAllConfig{patternConfigSuccess, patternConfigFail}
	configAllFail := &PatternAllConfig{patternConfigFail, patternConfigFail}
	configFiltersSuccess := &FiltersConfig{patternConfigSuccess, patternConfigSuccess}
	configFiltersSuccessAndFail := &FiltersConfig{patternConfigSuccess, patternConfigFail}
	configFiltersFail := &FiltersConfig{patternConfigFail, patternConfigFail}
	configWhitelistSuccess := &WhitelistConfig{patternConfigSuccess, patternConfigSuccess}
	configWhitelistSuccessAndFail := &WhitelistConfig{patternConfigSuccess, patternConfigFail}
	configWhitelistFail := &WhitelistConfig{patternConfigFail, patternConfigFail}

	r := createSuccessRequest()

	if !MustMatch(r, MustCreatePattern(configAllSuccess)) {
		t.Fatal(configAllSuccess)
	}
	if MustMatch(r, MustCreatePattern(configAllSuccessAndFail)) {
		t.Fatal(configAllSuccessAndFail)
	}
	if MustMatch(r, MustCreatePattern(configAllFail)) {
		t.Fatal(configAllFail)
	}
	if !MustMatch(r, &Filters{}) {
		t.Fatal()
	}
	if !MustMatch(r, MustCreatePattern(configFiltersSuccess)) {
		t.Fatal(configFiltersSuccess)
	}
	if !MustMatch(r, MustCreatePattern(configFiltersSuccessAndFail)) {
		t.Fatal(configFiltersSuccessAndFail)
	}
	if MustMatch(r, MustCreatePattern(configFiltersFail)) {
		t.Fatal(configFiltersFail)
	}
	if MustMatch(r, &Whitelist{}) {
		t.Fatal()
	}
	if !MustMatch(r, MustCreatePattern(configWhitelistSuccess)) {
		t.Fatal(configWhitelistSuccess)
	}
	if !MustMatch(r, MustCreatePattern(configWhitelistSuccessAndFail)) {
		t.Fatal(configWhitelistSuccessAndFail)
	}
	if MustMatch(r, MustCreatePattern(configWhitelistFail)) {
		t.Fatal(configWhitelistFail)
	}

}
