package persist

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

type Store interface {
	Start() error
	Stop() error
	SaveBytes(key string, data []byte) error
	LoadBytes(key string) ([]byte, error)
}

const FolderStoreExt = ".persist"
const FolderStoreFileMode = 0600
const FolderStoreMode = 0700

type FolderStore string

func (s FolderStore) Start() error {
	_, err := os.Stat(string(s))
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(string(s), FolderStoreMode)
		}
		return err
	}
	return nil
}
func (s FolderStore) Stop() error {
	return nil
}
func (s FolderStore) SaveBytes(key string, data []byte) error {
	name := filepath.Join(string(s), url.QueryEscape(key)+FolderStoreExt)
	return ioutil.WriteFile(name, data, FolderStoreFileMode)
}
func (s FolderStore) LoadBytes(key string) ([]byte, error) {
	name := filepath.Join(string(s), url.QueryEscape(key)+FolderStoreExt)
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return ioutil.ReadFile(name)
}
