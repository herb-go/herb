package cache_test

import (
	"testing"
	"time"
)

var testLaterLoader = func(key string) (interface{}, error) {
	time.Sleep(100 * time.Millisecond)
	return key, nil
}

func TestLaterLoader(t *testing.T) {
	var result string
	var result2 string
	var err error
	var err2 error
	c := newTestCache(3600)
	result = ""
	result2 = ""
	go func() {
		err = c.Load("test", &result, 0, testLaterLoader)

	}()
	go func() {
		err2 = c.Load("test", &result2, 0, testLaterLoader)

	}()
	time.Sleep(200 * time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if result != "test" {
		t.Fatal(result)
	}
	if err2 != nil {
		t.Fatal(err)
	}
	if result2 != "test" {
		t.Fatal(result)
	}
}
