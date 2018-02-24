package resource

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"
)

func LoadJSON(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
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
		return err
	}
	return nil
}
func MustLoadJSON(path string, v interface{}) {
	err := LoadJSON(path, v)
	if err != nil {
		panic(err)
	}
}
