package file

import (
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

type Assets struct {
	URLPrefix string
	Location  string
}

func (f *Assets) Save(id string, reader io.Reader) (string, error) {
	outfile := path.Join(f.Location, id)
	if !strings.HasPrefix(outfile, f.Location) {
		return "", errors.New("file local stote:unavailable file path")
	}
	dir := path.Dir(outfile)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0700)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	file, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE, 0700)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, reader)
	if err != nil {
		return "", err
	}
	return f.URL(id)
}
func (f *Assets) Load(id string, writer io.Writer) error {
	infile := path.Join(f.Location, id)
	if !strings.HasPrefix(infile, f.Location) {
		return errors.New("file local stote:unavailable file path")
	}

	file, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(writer, file)
	return err
}
func (f *Assets) Remove(id string) error {
	infile := path.Join(f.Location, id)
	if !strings.HasPrefix(infile, f.Location) {
		return errors.New("file local stote:unavailable file path")
	}
	return os.Remove(infile)

}
func (f *Assets) URL(id string) (string, error) {
	return path.Join(f.URLPrefix, id), nil
}
