package simplehttpserver

import (
	"net/http"
	"os"
	"path"
)

//ServeFile serve given file as http resopnse.
func ServeFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}

//ServeFolder serve  file in root path as http resopnse.
//Request folder with index.html will serve as index.html.
//Request folder without index.html will raise a http forbidden error(403).
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
