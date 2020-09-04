package persist

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestStore(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	f := FolderStore(filepath.Join(dir, "test"))
	err = f.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer f.Stop()
	d, err := f.LoadBytes("test")
	if d != nil || err != ErrNotFound {
		t.Fatal(d, err)
	}
	err = f.SaveBytes("test", []byte("testvalue"))
	if err != nil {
		t.Fatal(err)
	}
	d, err = f.LoadBytes("test")
	if string(d) != "testvalue" || err != nil {
		t.Fatal(d, err)
	}
	f = FolderStore(filepath.Join(dir, "test"))
	err = f.Start()
	if err != nil {
		t.Fatal(err)
	}
}
