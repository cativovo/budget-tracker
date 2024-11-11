package vite_test

import (
	"os"
	"testing"

	"github.com/cativovo/budget-tracker/internal/vite"
	"github.com/stretchr/testify/assert"
)

func TestAssets(t *testing.T) {
	testCases := []struct {
		want   string
		config vite.ViteConfig
	}{
		{
			config: vite.ViteConfig{
				DistFS: os.DirFS("testdata"),
				IsDev:  true,
			},
			want: `<script type="module" src="http://localhost:5173/@vite/client"></script><script type="module" src="http://localhost:5173/js/index.js"></script>`,
		},
		{
			config: vite.ViteConfig{
				DistFS: os.DirFS("testdata"),
				IsDev:  false,
			},
			want: `<link rel="stylesheet" href="/assets/index-Ckppcp00.css"><script type="module" src="/assets/index-CUUi8ibQ.js"></script>`,
		},
		{
			config: vite.ViteConfig{
				DistFS:   os.DirFS("testdata"),
				Manifest: "customdist/.vite/manifest.json",
				Assets:   "customdist/assets",
				IsDev:    false,
			},
			want: `<link rel="stylesheet" href="/assets/index-Ckppcp00.css"><script type="module" src="/assets/index-CUUi8ibQ.js"></script>`,
		},
	}

	for _, testCase := range testCases {
		v := vite.NewVite(testCase.config)
		assert.Equal(t, testCase.want, v.Assets(), "testProduction failed")
	}
}
