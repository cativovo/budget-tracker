package ui

type AssetsStore interface {
	Assets() string // inject using @templ.Raw
}
