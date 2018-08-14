package store

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

func (f *Assets) save(filename string, reader io.Reader, flag int) (string, int64, error) {
	outfile := path.Join(f.Location, filename)
	if !strings.HasPrefix(outfile, f.Location) {
		return "", 0, errors.New("file local stote:unavailable file path")
	}
	dir := path.Dir(outfile)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0700)
			if err != nil {
				return "", 0, err
			}
		} else {
			return "", 0, err
		}
	}
	file, err := os.OpenFile(outfile, flag, 0700)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()
	size, err := io.Copy(file, reader)
	if err != nil {
		return "", 0, err
	}
	return filename, size, nil
}

func (f *Assets) Save(filename string, reader io.Reader) (string, int64, error) {
	return f.save(filename, reader, os.O_WRONLY)
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

type AssetsStoreConfig struct {
	URLHost   string
	URLPrefix string
	Absolute  bool
	Root      string
	Location  string
}

func (c *AssetsStoreConfig) Create() (Driver, error) {
	driver := &Assets{}
	driver.URLPrefix = path.Join(c.URLHost, c.URLPrefix)
	if c.Absolute {
		driver.Location = path.Join("/", c.Root, c.Location)
	} else {
		root, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		driver.Location = path.Join(root, c.Root, c.Location)
	}
	_, err := os.Stat(driver.Location)
	if err != nil {
		return nil, err
	}
	return driver, nil
}
func init() {
	Register("assets", func(conf Config, prefix string) (Driver, error) {
		var err error
		c := &AssetsStoreConfig{}
		err = conf.Get(prefix+"URLHost", &c.URLHost)
		if err != nil {
			return nil, err
		}
		err = conf.Get(prefix+"URLPrefix", &c.URLPrefix)
		if err != nil {
			return nil, err
		}
		err = conf.Get(prefix+"Absolute", &c.Absolute)
		if err != nil {
			return nil, err
		}
		err = conf.Get(prefix+"URLHRootost", &c.Root)
		if err != nil {
			return nil, err
		}
		err = conf.Get(prefix+"Location", &c.Location)
		if err != nil {
			return nil, err
		}
		return c.Create()
	})
}
