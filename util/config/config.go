package config

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
)

func LoadJSON(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return NewError(path, err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var bytes = []byte{}
	err = nil
	var line string
	for err != io.EOF {
		line, err = r.ReadString(10)
		line = strings.TrimSpace(line)
		if len(line) > 2 && line[0:2] == "//" {
			continue
		}
		bytes = append(bytes, []byte(line)...)
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		return NewError(path, err)
	}
	return nil
}
func MustLoadJSON(path string, v interface{}) {
	err := LoadJSON(path, v)
	if err != nil {
		panic(err)
	}
}
func LoadTOML(path string, v interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return NewError(path, err)
	}
	err = toml.Unmarshal(data, v)
	if err != nil {
		return NewError(path, err)
	}
	return nil
}

func MustLoadTOML(path string, v interface{}) {
	err := LoadTOML(path, v)
	if err != nil {
		panic(err)
	}
}
