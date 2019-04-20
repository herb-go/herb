package store

import "testing"

func TestJSONConfig(t *testing.T) {
	c := NewOptionConfigJSON()
	c.Driver = "assets"
	err := c.Config.Set("Absolute", true)
	if err != nil {
		t.Fatal(err)
	}
	var result bool
	err = c.Config.Get("Absolute", &result)
	if err != nil {
		t.Fatal(err)
	}
	s := NewStore()
	err = s.Init(c)
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigMap(t *testing.T) {
	c := NewOptionConfigMap()
	c.Driver = "assets"
	err := c.Config.Set("Absolute", true)
	if err != nil {
		t.Fatal(err)
	}
	var result bool
	err = c.Config.Get("Absolute", &result)
	if err != nil {
		t.Fatal(err)
	}
	s := NewStore()
	err = s.Init(c)
	if err != nil {
		t.Fatal(err)
	}

}
