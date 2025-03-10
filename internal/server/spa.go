package server

import (
	"net/http"
	"os"
	"path"
	"strings"
)

func spaHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dist := os.DirFS(path.Join("ui", "dist"))
		f, err := dist.Open(strings.TrimPrefix(path.Clean(r.URL.Path), "/"))
		if err == nil {
			defer f.Close()
		}

		if os.IsNotExist(err) {
			r.URL.Path = "/"
		}

		http.FileServer(http.FS(dist)).ServeHTTP(w, r)
	}
}
