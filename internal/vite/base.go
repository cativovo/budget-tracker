package vite

import (
	"io/fs"
	"net/http"
	"os"

	budgettracker "github.com/cativovo/budget-tracker"
)

// https://github.com/olivere/vite
type Vite struct {
	assets   fs.FS
	manifest *Manifest
	IsDev    bool
}

func NewVite(isDev bool) Vite {
	var v Vite
	v.IsDev = isDev

	if !isDev {
		mf, err := budgettracker.Dist.Open("dist/.vite/manifest.json")
		if err != nil {
			panic(".vite/manifest.json not found: " + err.Error())
		}

		assets, err := fs.Sub(budgettracker.Dist, "dist/assets")
		if err != nil {
			panic(err)
		}

		manifest, err := ParseManifest(mf)
		if err != nil {
			panic("error ParseManifest: " + err.Error())
		}
		v.manifest = manifest
		v.assets = assets
	}

	return v
}

func (v Vite) Assets() string {
	if v.IsDev {
		return `
			<script type="module" src="http://localhost:5173/@vite/client"></script>
			<script type="module" src="http://localhost:5173/js/index.js"></script>
		`
	}

	chunk := v.manifest.GetEntryPoint()
	if chunk == nil {
		panic("entrypoint is nil")
	}

	css := v.manifest.GenerateCSS(chunk.Src)
	modules := v.manifest.GenerateModules(chunk.Src)
	preloadModules := v.manifest.GenerateCSS(chunk.Src)

	return css + modules + preloadModules
}

func (v Vite) AssetsFs() http.FileSystem {
	if v.IsDev {
		return http.FS(os.DirFS("js"))
	}

	return http.FS(v.assets)
}
