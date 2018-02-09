package misc

import (
	"net/http"
	"os"
	"path"
	"strings"
)

type fileFilesystem struct {
	Fs http.FileSystem
}

func (fs fileFilesystem) Open(name string) (http.File, error) {
	f, err := fs.Fs.Open(name)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if stat.IsDir() {
		err = f.Close()
		if err != nil {
			return nil, err
		}
		return fs.Open(path.Join(name, "index.html"))
	}

	return f, nil
}
func ServePrefixedFiles(prefix string, root http.FileSystem) http.HandlerFunc {
	fs := fileFilesystem{
		Fs: root,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		filepath := strings.TrimPrefix(url, prefix)

		if len(url) == len(filepath) {
			http.NotFound(w, r)
			return
		}
		r.URL.Path = filepath
		http.FileServer(fs).ServeHTTP(w, r)

	}
}
func ServeFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}
func ServeFolder(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file := path.Join(root, r.URL.Path)
		fi, err := os.Stat(file)
		if err == nil {
			if fi.IsDir() {
				ii, err := os.Stat(path.Join(file, "index.html"))
				if err != nil || ii.IsDir() {
					http.Error(w, http.StatusText(403), 403)
				}
			}
		}
		http.ServeFile(w, r, file)
	}
}
