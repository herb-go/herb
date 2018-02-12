package render

import "testing"

func TestData(t *testing.T) {
	data1 := Data{}
	data1.Set("test", "testvalue")
	if data1.Get("test").(string) != "testvalue" {
		t.Error(data1.Get("test"))
	}
	if data1.Get("test2") != nil {
		t.Error(data1.Get("test2"))
	}
	data1.Set("test2", "test2value")
	if data1.Get("test2").(string) != "test2value" {
		t.Error(data1.Get("test2"))
	}
	data1.Del("test2")
	if data1.Get("test2") != nil {
		t.Error(data1.Get("test2"))
	}
	data2 := Data{}
	data2.Set("test3", "test3value")
	data1.Merge(&data2)
	if data1.Get("test").(string) != "testvalue" {
		t.Error(data1.Get("test"))
	}
	if data1.Get("test3").(string) != "test3value" {
		t.Error(data1.Get("test3"))
	}
}

func TestNilData(t *testing.T) {
	data1 := new(Data)
	data1.Set("test", "testvalue")
	if data1.Get("test").(string) != "testvalue" {
		t.Error(data1.Get("test"))
	}
	data1 = new(Data)
	if data1.Get("test") != nil {
		t.Error(data1.Get("test"))
	}
	data1 = new(Data)
	data1.Del("test")
	if data1.Get("test") != nil {
		t.Error(data1.Get("test"))
	}
	data2 := new(Data)
	data2.Set("test2", "test2value")
	data1 = new(Data)
	data1.Merge(data2)
	if data1.Get("test2").(string) != "test2value" {
		t.Error(data1.Get("test2"))
	}
	data1 = new(Data)
	data2.Merge(data1)
	if data2.Get("test2").(string) != "test2value" {
		t.Error(data2.Get("test2"))
	}
}
