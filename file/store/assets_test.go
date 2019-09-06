package store

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestAssets(t *testing.T) {
	data := bytes.NewBuffer([]byte("filedata"))
	config := NewOptionConfigMap()
	config.Driver = "assets"
	config.Config.Set("URLHost", "http://www.test.com")
	config.Config.Set("URLPrefix", "test")
	config.Config.Set("Absolute", true)
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tmpdir)
	os.Mkdir(path.Join(tmpdir, "local"), 0700)
	config.Config.Set("Root", tmpdir)
	config.Config.Set("Location", "local")
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
	config := NewOptionConfigMap()
	config.Driver = "assets"
	config.Config.Set("URLHost", "http://www.test.com")
	config.Config.Set("URLPrefix", "test")
	config.Config.Set("Absolute", true)
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tmpdir)
	os.Mkdir(path.Join(tmpdir, "local"), 0700)
	config.Config.Set("Root", tmpdir)
	config.Config.Set("Location", "local")
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
	config := NewOptionConfigMap()
	config.Driver = "assets"
	config.Config.Set("URLHost", "http://www.test.com")
	config.Config.Set("URLPrefix", "test")
	config.Config.Set("Absolute", false)
	config.Config.Set("Root", "/")
	config.Config.Set("Location", "testdata")
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	driver, err := NewDriver(config.Driver, &config.Config, "")
	if err != nil {
		t.Fatal(err)
	}
	ad := driver.(*Assets)
	if ad.Location != path.Join(cwd, "testdata") {
		t.Fatal(ad.Location)
	}
}
func TestUrlEncode(t *testing.T) {
	config := NewOptionConfigMap()
	config.Driver = "assets"
	config.Config.Set("URLHost", "http://www.test.com")
	config.Config.Set("URLPrefix", "test")
	config.Config.Set("Absolute", true)
	store := New()
	err := store.Init(config)
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
