package server

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
)

var humaConfig = huma.DefaultConfig("Budget Tracker", "v0.0.1")

func init() {
	// TODO: setup properly
	// humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
	// 	"SessionAuth": {
	// 		Type: "apiKey",
	// 		In:   "cookie",
	// 		Name: SessionName,
	// 	},
	// }
	humaConfig.DocsPath = ""
	humaConfig.CreateHooks = nil
}

var security = []map[string][]string{
	{"SessionAuth": {}},
}

type humaHandler[I, O any] func(ctx context.Context, input *I) (*O, error)

func mountDocsRoute(r *chi.Mux) {
	r.Get("/docs", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.yaml"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`))
	})
}
