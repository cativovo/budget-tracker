package ui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var rootDir embed.FS

// DistDirFS contains the embedded dist directory files (without the "dist" prefix)
var DistDirFS fs.FS

func init() {
	d, err := fs.Sub(rootDir, "dist")
	if err != nil {
		panic(err)
	}

	DistDirFS = d
}
