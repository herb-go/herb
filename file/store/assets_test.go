package store

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestAssets(t *testing.T) {
	data := bytes.NewBuffer([]byte("filedata"))
	config := NewOptionConfig()
	config.Driver = "assets"
	conf := AssetsStoreConfig{
		URLHost:   "http://www.test.com",
		URLPrefix: "test",
		Absolute:  true,
	}
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tmpdir)
	os.Mkdir(path.Join(tmpdir, "local"), 0700)
	conf.Root = tmpdir
	conf.Location = "local"
	buf := bytes.NewBuffer(nil)
	err = json.NewEncoder(buf).Encode(conf)
	if err != nil {
		panic(err)
	}
	config.Config = json.NewDecoder(buf).Decode
	store := New()
	err = store.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	id, length, err := store.Save("test", data)
	if err != nil {
		t.Fatal(err)
	}

	if length != int64(len("filedata")) {
		t.Fatal(data)
	}
	if id == "" {
		t.Fatal(id)
	}

	_, length, err = store.Save("/test2/test", data)
	if err != nil {
		t.Fatal(err)
	}
	reader, err := store.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	reader.Close()
	if string(content) != "filedata" {
		t.Fatal(content)
	}
	storefile := store.File(id)
	url, err := storefile.URL()
	if err != nil {
		t.Fatal(err)
	}
	if url != "http:/www.test.com/test/test" {
		t.Fatal(url)
	}
	url, err = storefile.URL()
	if err != nil {
		t.Fatal(err)
	}
	if url != "http:/www.test.com/test/test" {
		t.Fatal(url)
	}

	err = store.Remove(id)
	if err != nil {
		t.Fatal(err)
	}
	err = store.Remove(id)
	if err == nil || err.(*Error) == nil || err.(*Error).Type != ErrorTypeNotExists || err.(*Error).File != id {
		t.Fatal(err)
	}
}

func TestFileNameFail(t *testing.T) {
	data := bytes.NewBuffer([]byte("filedata"))
	config := NewOptionConfig()
	config.Driver = "assets"
	conf := AssetsStoreConfig{
		URLHost:   "http://www.test.com",
		URLPrefix: "test",
		Absolute:  true,
	}

	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tmpdir)
	os.Mkdir(path.Join(tmpdir, "local"), 0700)
	conf.Root = tmpdir
	conf.Location = "local"
	buf := bytes.NewBuffer(nil)
	err = json.NewEncoder(buf).Encode(conf)
	if err != nil {
		panic(err)
	}
	config.Config = json.NewDecoder(buf).Decode

	store := New()
	err = store.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = store.Save("../test", data)
	if err == nil || err.(*Error) == nil || err.(*Error).Type != ErrorTypeUnavailableID || err.(*Error).File != "../test" {
		t.Fatal(err)
	}
	_, err = store.Load("../test")
	if err == nil || err.(*Error) == nil || err.(*Error).Type != ErrorTypeUnavailableID || err.(*Error).File != "../test" {
		t.Fatal(err)
	}
	err = store.Remove("../test")
	if err == nil || err.(*Error) == nil || err.(*Error).Type != ErrorTypeUnavailableID || err.(*Error).File != "../test" {
		t.Fatal(err)
	}
}

func TestAbsolute(t *testing.T) {
	config := NewOptionConfig()
	config.Driver = "assets"
	conf := AssetsStoreConfig{
		URLHost:   "http://www.test.com",
		URLPrefix: "test",
		Absolute:  false,
		Root:      "/",
		Location:  "testdata",
	}
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(conf)
	if err != nil {
		panic(err)
	}
	config.Config = json.NewDecoder(buf).Decode
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	driver, err := NewDriver(config.Driver, config.Config)
	if err != nil {
		t.Fatal(err)
	}
	ad := driver.(*Assets)
	if ad.Location != path.Join(cwd, "testdata") {
		t.Fatal(ad.Location)
	}
}
func TestUrlEncode(t *testing.T) {
	config := NewOptionConfig()
	config.Driver = "assets"
	conf := AssetsStoreConfig{
		URLHost:   "http://www.test.com",
		URLPrefix: "test",
		Absolute:  true,
	}
	store := New()
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(conf)
	if err != nil {
		panic(err)
	}
	config.Config = json.NewDecoder(buf).Decode
	err = store.Init(config)
	if err != nil {
		t.Fatal(err)
	}
	url, err := store.URL("%")
	if err != nil {
		t.Fatal(err)
	}
	if url != "http:/www.test.com/test/%25" {
		t.Fatal(url)
	}

	url, err = store.URL("abc")
	if err != nil {
		t.Fatal(err)
	}
	if url != "http:/www.test.com/test/abc" {
		t.Fatal(url)
	}
}
