package cache_test

import (
	"sync"
	"testing"
	"time"

	"github.com/herb-go/herb/cache"
)

var testLaterLoader = func(key string) (interface{}, error) {
	time.Sleep(100 * time.Millisecond)
	return key, nil
}

func TestCloneUtil(t *testing.T) {
	u := cache.NewUtil()
	uc := u.Clone()
	uc.NodeFactory = func(cache.Cacheable, string) *cache.Node {
		return nil
	}
	if u.NodeFactory != nil || uc.NodeFactory == nil {
		t.Fatal(u, uc)
	}
}
func TestLaterLoader(t *testing.T) {
	var result string
	var result2 string
	var result3 string
	var result4 string
	var result5 string

	var err error
	var err2 error
	var err3 error
	var err4 error
	var err5 error
	c := newTestCache(3600)
	result = ""
	result2 = ""
	result3 = ""
	result4 = ""
	result5 = ""
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = c.Load("test", &result, 0, testLaterLoader)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		err2 = c.Load("test", &result2, 0, testLaterLoader)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		err3 = c.Load("test2", &result3, 0, testLaterLoader)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		err4 = c.Load("test2", &result4, 0, testLaterLoader)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		err5 = c.Load("test3", &result5, 0, testLaterLoader)
		wg.Done()
	}()
	wg.Wait()
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
		t.Fatal(result2)
	}
	if err3 != nil {
		t.Fatal(err)
	}
	if result3 != "test2" {
		t.Fatal(result3)
	}
	if err4 != nil {
		t.Fatal(err)
	}
	if result4 != "test2" {
		t.Fatal(result4)
	}
	if err5 != nil {
		t.Fatal(err)
	}
	if result5 != "test3" {
		t.Fatal(result5)
	}
}
