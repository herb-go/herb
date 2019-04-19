package blocker

import "testing"
import "encoding/json"

func TestConfig(t *testing.T) {
	var config = `
	[{
		"StatusCode":403,
		"Limit":10,
		"DurationInSecond":60
	},{
		"StatusCode":-1,
		"Limit":100,
		"DurationInSecond":60
	}]`

	r := NewRules()
	err := json.Unmarshal([]byte(config), r)
	if err != nil {
		t.Fatal(err)
	}
	b := New(newTestCache(1 * 3600))
	if err != nil {
		t.Fatal(err)
	}
	err = r.ApplyTo(b)
	if len(b.config) != 2 {
		t.Fatal(b.config)
	}
	v, ok := b.config[403]
	if ok == false || v.max != 10 || v.ttlSecond != 60 {
		t.Fatal(b.config)
	}
	v, ok = b.config[-1]
	if ok == false || v.max != 100 || v.ttlSecond != 60 {
		t.Fatal(b.config)
	}
}
