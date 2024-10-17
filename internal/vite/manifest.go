package vite

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// https://github.com/olivere/vite/blob/main/manifest.go

// Manifest file as written by vite build, as described in the [Vite Manifest].
//
// It is required for backend integration as described in
// [Vite Backend Integration].
//
// [Vite Manifest]: https://vitejs.dev/guide/api-plugin.html#manifest
// [Vite Backend Integration]: https://vitejs.dev/guide/backend-integration.html
type Manifest map[string]*Chunk

// A Chunk is a single entry in the manifest.
type Chunk struct {
	File           string   `json:"file"`
	Name           string   `json:"name"`
	Src            string   `json:"src"`
	CSS            []string `json:"css"`
	Imports        []string `json:"imports"`
	DynamicImports []string `json:"dynamicImports"`
	IsDynamicEntry bool     `json:"isDynamicEntry"`
	IsEntry        bool     `json:"isEntry"`
}

// ParseManifest parses the manifest file.
func ParseManifest(r io.Reader) (*Manifest, error) {
	var m Manifest
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

// GetEntryPoint returns the entry point from the Vite manifest.
func (m Manifest) GetEntryPoint() *Chunk {
	for _, chunk := range m {
		if chunk.IsEntry {
			return chunk
		}
	}
	return nil
}

// GetEntryPoints returns the entry points from the manifest.
func (m Manifest) GetEntryPoints() []*Chunk {
	var entryPoints []*Chunk
	for _, chunk := range m {
		if chunk.IsEntry {
			entryPoints = append(entryPoints, chunk)
		}
	}
	return entryPoints
}

// GetChunk returns the chunk with the given name from the manifest.
//
// The name is the name of the source file.
func (m Manifest) GetChunk(name string) (*Chunk, bool) {
	chunk, ok := m[name]
	return chunk, ok
}

// GenerateCSS generates the CSS links for the given chunk.
//
// The name is the name of the source file, e.g. "src/main.tsx".
func (m Manifest) GenerateCSS(name string) string {
	var sb strings.Builder
	seen := make(map[string]bool)

	var addCSS func(string)
	addCSS = func(name string) {
		if seen[name] {
			return
		}
		seen[name] = true

		chunk, ok := m[name]
		if !ok {
			return
		}

		for _, css := range chunk.CSS {
			fmt.Fprintf(&sb, `<link rel="stylesheet" href="/%s">`, css)
		}

		for _, imp := range chunk.Imports {
			addCSS(imp)
		}
	}

	addCSS(name)

	return sb.String()
}

// GenerateModules generates the module scripts for the given chunk.
//
// The name is the name of the source file, e.g. "src/main.tsx".
func (m Manifest) GenerateModules(name string) string {
	chunk, ok := m[name]
	if !ok {
		return ""
	}

	var sb strings.Builder
	if chunk.File != "" {
		fmt.Fprintf(&sb, `<script type="module" src="/%s"></script>`, chunk.File)
	}

	return sb.String()
}

// GeneratePreloadModules generates the preload modules for the given chunk.
//
// The name is the name of the source file, e.g. "src/main.tsx".
func (m Manifest) GeneratePreloadModules(name string) string {
	var sb strings.Builder
	seen := make(map[string]bool)

	var addModulePreload func(string)
	addModulePreload = func(name string) {
		if seen[name] {
			return
		}
		seen[name] = true

		chunk, ok := m[name]
		if !ok {
			return
		}

		if chunk.File != "" {
			fmt.Fprintf(&sb, `<link rel="modulepreload" href="/%s">`, chunk.File)
		}

		for _, imp := range chunk.Imports {
			addModulePreload(imp)
		}
	}

	addModulePreload(name)

	return sb.String()
}
