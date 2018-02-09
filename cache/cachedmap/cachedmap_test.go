package cachedmap

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/herb-go/herb/cache"
	_ "github.com/herb-go/herb/cache/drivers/freecache"
)

func newTestCache(ttl int64) *cache.Cache {
	config := json.RawMessage("{\"Size\": 10000000}")
	c := cache.New()
	err := c.Open("freecache", config, ttl)
	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	return c
}

type testmodel struct {
	Keyword string
	Content int
}

const valueKey = "valueKey"
const valueKeyAadditional = "valueKeyAadditional"
const valueKeyChanged = "valueKeyChanged"
const wrongDataKey = "wrongdata"

var WrongData = []byte("wrongdata")

const startValue = 1
const changedValue = 2
const mapCreatorKeyword = "mapCreatorKeyword"

const creatorKeyword = "creatorKeyword"

var rawData map[string]int

type testmodelmap map[string]*testmodel

func (m *testmodelmap) NewMapElement(key string) error {
	(*m)[key] = &testmodel{
		Keyword: mapCreatorKeyword,
		Content: 0,
	}
	return nil
}
func (m *testmodelmap) LoadMapElements(keys ...string) error {
	for _, v := range keys {
		(*m)[v].Content = rawData[v]
	}
	return nil
}
func creator(m *testmodelmap) func(key string) error {
	return func(key string) error {
		(*m)[key] = &testmodel{
			Keyword: creatorKeyword,
			Content: 0,
		}
		return nil
	}
}

func loader(m *testmodelmap) func(keys ...string) error {
	return func(keys ...string) error {
		for _, v := range keys {
			(*m)[v].Content = rawData[v]
		}
		return nil
	}
}
func TestMap(t *testing.T) {
	rawData = map[string]int{
		valueKey:            startValue,
		valueKeyAadditional: startValue,
		valueKeyChanged:     startValue,
	}
	c := newTestCache(-1)
	var err error
	var tm = testmodelmap{}
	err = LoadCachedMap(&tm, c, valueKey, valueKeyAadditional, valueKeyAadditional)
	if err != nil {
		t.Fatal(err)
	}

	if val := tm[valueKey].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKey].Keyword; val != mapCreatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Keyword; val != mapCreatorKeyword {
		t.Error(val)
	}
	delete(tm, valueKeyAadditional)
	rawData[valueKey] = changedValue
	rawData[valueKeyAadditional] = changedValue
	rawData[valueKeyChanged] = changedValue
	err = LoadCachedMap(&tm, c, valueKey, valueKeyAadditional, valueKeyChanged)
	if err != nil {
		t.Fatal(err)
	}
	if val := tm[valueKey].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKey].Keyword; val != mapCreatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Keyword; val != mapCreatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyChanged].Content; val != changedValue {
		t.Error(val)
	}
	if val := tm[valueKeyChanged].Keyword; val != mapCreatorKeyword {
		t.Error(val)
	}
	c.SetBytesValue(wrongDataKey, WrongData, 0)
	err = LoadCachedMap(&tm, c, wrongDataKey)
	if err == nil {
		t.Fatal(err)
	}

}

func TestMapLoad(t *testing.T) {
	rawData = map[string]int{
		valueKey:            startValue,
		valueKeyAadditional: startValue,
		valueKeyChanged:     startValue,
	}
	c := newTestCache(-1)
	var err error
	var tm = testmodelmap{}
	err = Load(&tm, c, loader(&tm), creator(&tm), valueKey, valueKeyAadditional)
	if err != nil {
		t.Fatal(err)
	}
	if val := tm[valueKey].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKey].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	rawData[valueKey] = changedValue
	rawData[valueKeyAadditional] = changedValue
	rawData[valueKeyChanged] = changedValue
	err = Load(&tm, c, loader(&tm), creator(&tm), valueKeyAadditional, valueKeyChanged)
	if err != nil {
		t.Fatal(err)
	}
	if val := tm[valueKey].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKey].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyChanged].Content; val != changedValue {
		t.Error(val)
	}
	if val := tm[valueKeyChanged].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	var tm2 = testmodelmap{}
	err = Load(&tm2, c, loader(&tm2), creator(&tm2), valueKeyAadditional, valueKeyChanged)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMapNodeLoad(t *testing.T) {
	rawData = map[string]int{
		valueKey:            startValue,
		valueKeyAadditional: startValue,
		valueKeyChanged:     startValue,
	}
	c := newTestCache(-1).Node("test")
	var err error
	var tm = testmodelmap{}
	err = Load(&tm, c, loader(&tm), creator(&tm), valueKey, valueKeyAadditional)
	if err != nil {
		t.Fatal(err)
	}

	if val := tm[valueKey].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKey].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	rawData[valueKey] = changedValue
	rawData[valueKeyAadditional] = changedValue
	rawData[valueKeyChanged] = changedValue
	err = Load(&tm, c, loader(&tm), creator(&tm), valueKeyAadditional, valueKeyChanged)
	if err != nil {
		t.Fatal(err)
	}
	if val := tm[valueKey].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKey].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyChanged].Content; val != changedValue {
		t.Error(val)
	}
	if val := tm[valueKeyChanged].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	var tm2 = testmodelmap{}
	err = Load(&tm2, c, loader(&tm2), creator(&tm2), valueKeyAadditional, valueKeyChanged)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMapCollectionLoad(t *testing.T) {
	rawData = map[string]int{
		valueKey:            startValue,
		valueKeyAadditional: startValue,
		valueKeyChanged:     startValue,
	}
	c := newTestCache(3600).Collection("test")
	var err error
	var tm = testmodelmap{}
	err = Load(&tm, c, loader(&tm), creator(&tm), valueKey, valueKeyAadditional)
	if err != nil {
		t.Fatal(err)
	}

	if val := tm[valueKey].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKey].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	rawData[valueKey] = changedValue
	rawData[valueKeyAadditional] = changedValue
	rawData[valueKeyChanged] = changedValue
	err = Load(&tm, c, loader(&tm), creator(&tm), valueKeyAadditional, valueKeyChanged)
	if err != nil {
		t.Fatal(err)
	}
	if val := tm[valueKey].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKey].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyChanged].Content; val != changedValue {
		t.Error(val)
	}
	if val := tm[valueKeyChanged].Keyword; val != creatorKeyword {
		t.Error(val)
	}
}

var rawString map[string]string

type teststringmap map[string]string

func (m *teststringmap) NewMapElement(key string) error {
	(*m)[key] = ""
	return nil
}
func (m *teststringmap) LoadMapElements(keys ...string) error {
	for _, v := range keys {
		(*m)[v] = rawString[v]
	}
	return nil
}
func TestString(t *testing.T) {
	var emptyKey = "empty"
	rawString = map[string]string{
		emptyKey:            "",
		valueKey:            valueKey,
		valueKeyAadditional: valueKeyAadditional,
		valueKeyChanged:     valueKeyChanged,
	}
	c := newTestCache(3600).Collection("test")
	var err error
	var tm = teststringmap{}
	err = LoadCachedMap(&tm, c, valueKey, valueKeyAadditional, emptyKey, valueKeyChanged)
	if err != nil {
		t.Fatal(err)
	}

	if val := tm[valueKey]; val != valueKey {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional]; val != valueKeyAadditional {
		t.Error(val)
	}
	if val := tm[emptyKey]; val != "" {
		t.Error(val)
	}

	err = LoadCachedMap(&tm, c, valueKey, valueKeyAadditional)
	if err != nil {
		t.Fatal(err)
	}

	if val := tm[valueKey]; val != valueKey {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional]; val != valueKeyAadditional {
		t.Error(val)
	}
	var tm2 = teststringmap{}
	err = LoadCachedMap(&tm2, c, valueKey, valueKeyAadditional, valueKeyChanged)
	if err != nil {
		t.Fatal(err)
	}

}

var rawBytes map[string][]byte

type testbytesmap map[string][]byte

func (m *testbytesmap) NewMapElement(key string) error {
	(*m)[key] = []byte{}
	return nil
}
func (m *testbytesmap) LoadMapElements(keys ...string) error {
	for _, v := range keys {
		(*m)[v] = rawBytes[v]
	}
	return nil
}
func TestBytes(t *testing.T) {
	rawBytes = map[string][]byte{
		valueKey:            []byte(valueKey),
		valueKeyAadditional: []byte(valueKeyAadditional),
		valueKeyChanged:     []byte(valueKeyChanged),
	}
	c := newTestCache(3600).Collection("test")
	var err error
	var tm = testbytesmap{}
	err = LoadCachedMap(&tm, c, valueKey, valueKeyAadditional)
	if err != nil {
		t.Fatal(err)
	}

	if val := tm[valueKey]; bytes.Compare(val, []byte(valueKey)) != 0 {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional]; bytes.Compare(val, []byte(valueKeyAadditional)) != 0 {
		t.Error(val)
	}
	err = LoadCachedMap(&tm, c, valueKey, valueKeyAadditional)
	if err != nil {
		t.Fatal(err)
	}

	if val := tm[valueKey]; bytes.Compare(val, []byte(valueKey)) != 0 {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional]; bytes.Compare(val, []byte(valueKeyAadditional)) != 0 {
		t.Error(val)
	}
	var tm2 = testbytesmap{}
	err = LoadCachedMap(&tm2, c, valueKey, valueKeyAadditional)
	if err != nil {
		t.Fatal(err)
	}
}
