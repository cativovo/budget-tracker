package vite

import (
	"io/fs"
	"net/http"
	"os"
	"strings"
)

type ViteConfig struct {
	DistFS   fs.FS
	Manifest string
	Assets   string
	IsDev    bool
}

// https://github.com/olivere/vite
type Vite struct {
	assets   fs.FS
	manifest *Manifest
	isDev    bool
}

func NewVite(config ViteConfig) Vite {
	var v Vite
	v.isDev = config.IsDev

	if !v.isDev {
		if config.Manifest == "" {
			config.Manifest = "dist/.vite/manifest.json"
		}

		if config.Assets == "" {
			config.Assets = "dist/assets"
		}

		manifestFile, err := config.DistFS.Open(config.Manifest)
		if err != nil {
			panic("error opening manifest file: " + err.Error())
		}
		defer manifestFile.Close()

		manifest, err := ParseManifest(manifestFile)
		if err != nil {
			panic("error parsing manifest: " + err.Error())
		}
		v.manifest = manifest

		assets, err := fs.Sub(config.DistFS, config.Assets)
		if err != nil {
			panic("error getting assets: " + err.Error())
		}
		v.assets = assets
	}

	return v
}

func (v Vite) Assets() string {
	if v.isDev {
		return `<script type="module" src="http://localhost:5173/@vite/client"></script><script type="module" src="http://localhost:5173/js/index.js"></script>`
	}

	chunk := v.manifest.GetEntryPoint()
	if chunk == nil {
		panic("manifest entrypoint is nil")
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
	if v.isDev {
		return http.FS(os.DirFS("js"))
	}

	return http.FS(v.assets)
}
