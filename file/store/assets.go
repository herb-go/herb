package store

import (
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

//Assets local file store.
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

//Save save data form reader to named file.
//Return file id ,file size and any error if raised.
func (f *Assets) Save(filename string, reader io.Reader) (string, int64, error) {
	return f.save(filename, reader, os.O_WRONLY|os.O_CREATE)
}

//Load load file with given id and write to writer.
//Return any error if raised.
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

//Remove remove file by id.
//Return any error if raised.
func (f *Assets) Remove(id string) error {
	infile := path.Join(f.Location, id)
	if !strings.HasPrefix(infile, f.Location) {
		return errors.New("file local stote:unavailable file path")
	}
	return os.Remove(infile)

}

//URL convert file id to file url.
//Return file url and any error if raised.
func (f *Assets) URL(id string) (string, error) {
	return path.Join(f.URLPrefix, id), nil
}

//AssetsStoreConfig local file store config
type AssetsStoreConfig struct {
	//URLHost file url host.
	URLHost string
	//URLPrefix file url path prefix
	URLPrefix string
	//Absolute if filepath is Absolute.
	Absolute bool
	//Root file root path.
	//if Absolute is true,file root is based from root path.
	//Otherwie file root is based from current working direction.
	Root string
	//Location file sub  folder which stored in.
	Location string
}

//Create create new local file driver.
//Return created driver and any error if raised.
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
		err = conf.Get(prefix+"Root", &c.Root)
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
