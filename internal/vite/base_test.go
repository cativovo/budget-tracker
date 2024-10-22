package vite_test

import (
	"os"
	"testing"

	"github.com/cativovo/budget-tracker/internal/vite"
	"github.com/stretchr/testify/assert"
)

func TestAssets(t *testing.T) {
	viteManifest, err := os.Open("testdata/dist/.vite/manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	defer viteManifest.Close()
	assets := os.DirFS("testdata/dist/assets")

	testDev := func() {
		v := vite.NewVite(true, viteManifest, assets)
		want := `
			<script type="module" src="http://localhost:5173/@vite/client"></script>
			<script type="module" src="http://localhost:5173/js/index.js"></script>
		`
		assert.Equal(t, want, v.Assets(), "testDev failed")
	}

	testProduction := func() {
		v := vite.NewVite(false, viteManifest, assets)
		want := `<link rel="stylesheet" href="/assets/index-Ckppcp00.css"><script type="module" src="/assets/index-CUUi8ibQ.js"></script>`
		assert.Equal(t, want, v.Assets(), "testProduction failed")
	}

	testDev()
	testProduction()
}
