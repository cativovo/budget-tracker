package vite_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cativovo/budget-tracker/internal/vite"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCSS(t *testing.T) {
	testCases := []struct {
		manifestFile string
		want         string
		sourceFile   string
	}{
		{
			manifestFile: "manifest1.json",
			want:         `<link rel="stylesheet" href="/assets/index-Ckppcp00.css">`,
			sourceFile:   "js/index.js",
		},
		{
			manifestFile: "manifest2.json",
			want:         `<link rel="stylesheet" href="/assets/foo-5UjPuW-k.css"><link rel="stylesheet" href="/assets/shared-ChJ_j-JJ.css">`,
			sourceFile:   "views/foo.js",
		},
		{
			manifestFile: "manifest2.json",
			want:         `<link rel="stylesheet" href="/assets/shared-ChJ_j-JJ.css">`,
			sourceFile:   "views/bar.js",
		},
	}

	for _, testCase := range testCases {
		manifest := parseManifest(t, testCase.manifestFile)
		assert.Equal(t, testCase.want, manifest.GenerateCSS(testCase.sourceFile))
	}
}

func TestGenerateGenerateModules(t *testing.T) {
	testCases := []struct {
		manifestFile string
		want         string
		sourceFile   string
	}{
		{
			manifestFile: "manifest1.json",
			want:         `<script type="module" src="/assets/index-CUUi8ibQ.js"></script>`,
			sourceFile:   "js/index.js",
		},
		{
			manifestFile: "manifest2.json",
			want:         `<script type="module" src="/assets/foo-BRBmoGS9.js"></script>`,
			sourceFile:   "views/foo.js",
		},
		{
			manifestFile: "manifest2.json",
			want:         `<script type="module" src="/assets/bar-gkvgaI9m.js"></script>`,
			sourceFile:   "views/bar.js",
		},
		{
			manifestFile: "manifest2.json",
			want:         `<script type="module" src="/assets/baz-B2H3sXNv.js"></script>`,
			sourceFile:   "baz.js",
		},
	}

	for _, testCase := range testCases {
		manifest := parseManifest(t, testCase.manifestFile)
		assert.Equal(t, testCase.want, manifest.GenerateModules(testCase.sourceFile))
	}
}

func TestGenerateGeneratePreloadModules(t *testing.T) {
	testCases := []struct {
		manifestFile string
		want         string
		sourceFile   string
	}{
		{
			manifestFile: "manifest1.json",
			want:         `<link rel="modulepreload" href="/assets/index-CUUi8ibQ.js">`,
			sourceFile:   "js/index.js",
		},
		{
			manifestFile: "manifest2.json",
			want:         `<link rel="modulepreload" href="/assets/foo-BRBmoGS9.js"><link rel="modulepreload" href="/assets/shared-B7PI925R.js">`,
			sourceFile:   "views/foo.js",
		},
		{
			manifestFile: "manifest2.json",
			want:         `<link rel="modulepreload" href="/assets/bar-gkvgaI9m.js"><link rel="modulepreload" href="/assets/shared-B7PI925R.js">`,
			sourceFile:   "views/bar.js",
		},
		{
			manifestFile: "manifest2.json",
			want:         `<link rel="modulepreload" href="/assets/baz-B2H3sXNv.js">`,
			sourceFile:   "baz.js",
		},
	}

	for _, testCase := range testCases {
		manifest := parseManifest(t, testCase.manifestFile)
		assert.Equal(t, testCase.want, manifest.GeneratePreloadModules(testCase.sourceFile))
	}
}

func parseManifest(t *testing.T, filename string) *vite.Manifest {
	mf, err := os.Open(fmt.Sprintf("testdata/%s", filename))
	if err != nil {
		t.Fatal(err)
	}
	defer mf.Close()

	manifest, err := vite.ParseManifest(mf)
	if err != nil {
		t.Fatal(err)
	}

	return manifest
}
