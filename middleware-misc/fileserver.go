package misc

import (
	"net/http"
	"os"
	"strings"
)

type justFilesFilesystem struct {
	Fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
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
		return nil, os.ErrNotExist
	}

	return f, nil
}
func ServePrefixedFiles(prefix string, root http.FileSystem) http.HandlerFunc {
	fs := justFilesFilesystem{
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
func ServeFolder(root http.FileSystem) http.HandlerFunc {
	fs := justFilesFilesystem{
		Fs: root,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(fs).ServeHTTP(w, r)
	}
}
