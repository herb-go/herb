package cache

import "testing"
import "encoding/json"
import "time"

func TestDummyCache(t *testing.T) {
	c := New()
	testKey := "key"
	testData := []string{"123"}
	testBytes := []byte("123")
	testTTL := time.Duration(-1)
	testIncrement := int64(1)
	testIntValue := int64(2)
	var model string
	err := c.Open("dummycache", json.RawMessage(""), int64(testTTL))
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.SearchByPrefix(testKey)
	if err != ErrFeatureNotSupported {
		t.Errorf("Search by prefeix error %s", err)
	}
	err = c.Set(testKey, testData, testTTL)
	if err != nil {
		t.Errorf("Set error %s", err)
	}
	err = c.Update(testKey, testData, testTTL)
	if err != nil {
		t.Errorf("Update error %s", err)
	}
	err = c.Get(testKey, &model)
	if err != ErrNotFound {
		t.Errorf("Get error %s", err)
	}
	err = c.SetBytesValue(testKey, testBytes, testTTL)
	if err != nil {
		t.Errorf("SetBytesValue error %s", err)
	}
	err = c.UpdateBytesValue(testKey, testBytes, testTTL)
	if err != nil {
		t.Errorf("UpdateBytesValue error %s", err)
	}
	_, err = c.GetBytesValue(testKey)
	if err != ErrNotFound {
		t.Errorf("GetBytesValue error %s", err)
	}
	err = c.Del(testKey)
	if err != nil {
		t.Errorf("Del error %s", err)
	}
	_, err = c.IncrCounter(testKey, testIncrement, testTTL)
	if err != nil {
		t.Errorf("IncrCounter error %s", err)
	}
	err = c.SetCounter(testKey, testIntValue, testTTL)
	if err != nil {
		t.Errorf("SetCounter error %s", err)
	}
	_, err = c.GetCounter(testKey)
	if err != ErrNotFound {
		t.Errorf("GetCounter error %s", err)
	}
	err = c.DelCounter(testKey)
	if err != nil {
		t.Errorf("DelCounter error %s", err)
	}
	err = c.Flush()
	if err != nil {
		t.Errorf("Flush error %s", err)
	}
	err = c.Close()
	if err != nil {
		t.Errorf("Close error %s", err)
	}
}
