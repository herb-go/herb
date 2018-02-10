package misc

import (
	"net/http"
	"os"
	"path"
)

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
