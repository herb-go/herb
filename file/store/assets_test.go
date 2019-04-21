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
	config := NewOptionConfigJSON()
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
}
