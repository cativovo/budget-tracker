package vite

import (
	"io/fs"
	"net/http"
	"os"
	"strings"
)

// https://github.com/olivere/vite
type Vite struct {
	assets   fs.FS
	manifest *Manifest
	IsDev    bool
}

func NewVite(isDev bool, mf fs.File, assetsFS fs.FS) Vite {
	var v Vite
	v.IsDev = isDev

	if !isDev {
		manifest, err := ParseManifest(mf)
		if err != nil {
			panic("error ParseManifest: " + err.Error())
		}
		v.manifest = manifest
		v.assets = assetsFS
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

	var preloadModules strings.Builder
	for _, name := range chunk.Imports {
		preloadModules.WriteString(v.manifest.GeneratePreloadModules(name))
	}

	return css + modules + preloadModules.String()
}

func (v Vite) AssetsFs() http.FileSystem {
	if v.IsDev {
		return http.FS(os.DirFS("js"))
	}

	return http.FS(v.assets)
}
