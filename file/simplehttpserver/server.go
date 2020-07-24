package simplehttpserver

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func isSlashRune(r rune) bool { return r == '/' || r == '\\' }

func ContainsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}
	return false
}
func renderError(w http.ResponseWriter, err error) {

	if os.IsNotExist(err) {
		http.Error(w, http.StatusText(404), 404)
	} else if os.IsPermission(err) {
		http.Error(w, http.StatusText(404), 403)
	} else {
		panic(err)
	}
}
func Download(w http.ResponseWriter, r *http.Request, path string, name string) {
	if ContainsDotDot(r.URL.Path) {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	stat, err := os.Stat(path)
	if err != nil {
		renderError(w, err)
		return
	}
	if stat.IsDir() {
		http.Error(w, http.StatusText(404), 403)
		return
	}
	file, err := os.Open(path)
	if err != nil {
		renderError(w, err)
		return
	}
	defer file.Close()
	if name == "" {
		name = filepath.Base(path)
	}
	http.ServeContent(w, r, name, stat.ModTime(), file)
}

//ServeFile serve given file as http resopnse.
func ServeFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Download(w, r, path, "")
	}
}

//ServeFolder serve files in root path as http resopnse.
//Request folder with index.html will serve as index.html.
//Request folder without index.html will raise a http forbidden error(403).
func ServeFolder(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file := path.Join(root, r.URL.Path)
		stat, err := os.Stat(file)
		if err != nil {
			renderError(w, err)
			return
		}
		if stat.IsDir() {
			filestat, err := os.Stat(path.Join(file, "index.html"))
			if err == nil && !filestat.IsDir() {
				file = path.Join(file, "index.html")
			}
		}
		Download(w, r, file, "")
	}
}
