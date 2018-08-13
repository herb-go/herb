package file

import (
	"io"
)

type File struct {
	Store Store
	ID    string
	url   string
}

func (f *File) SetURL(url string) {
	f.url = url
}
func (f *File) URL() (url string, err error) {
	if f.url != "" {
		return f.url, nil
	}
	return f.Store.URL(f.ID)
}

type Store interface {
	Save(id string, reader io.Reader) (string, error)
	Load(id string, reader io.Writer) error
	Remove(id string) error
	URL(id string) (string, error)
}
