package server

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/cativovo/budget-tracker/ui"
)

func spaHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := ui.DistDirFS.Open(strings.TrimPrefix(path.Clean(r.URL.Path), "/"))
		if err == nil {
			defer f.Close()
		}

		if os.IsNotExist(err) {
			r.URL.Path = "/"
		}

		http.FileServer(http.FS(ui.DistDirFS)).ServeHTTP(w, r)
	}
}
