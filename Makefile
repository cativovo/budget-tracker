ENV_FILE ?= .env.development.local
include $(ENV_FILE)
export

seed:
	go run ./cmd/seed

cleandb:
	go run ./cmd/seed -c

liveassets:
	pnpm run dev

livetempl:
	templ generate --watch

liveserver:
	go run github.com/air-verse/air@v1.60.0 \
	--build.cmd "go build -o ./tmp/bin/main ./cmd/app" \
	--build.bin "./tmp/bin/main" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

live:
	make -j3 liveassets livetempl liveserver

generate:
	go run github.com/go-jet/jet/v2/cmd/jet@v2.12.0 -dsn=postgresql://$$POSTGRES_USER:$$POSTGRES_PASSWORD@$$POSTGRES_HOST:$$POSTGRES_PORT/$$POSTGRES_DB?sslmode=$$POSTGRES_SSL -path=internal/repository
	go run github.com/a-h/templ/cmd/templ@v0.2.793 generate

build: generate
	pnpm run build
	go build -o app ./cmd/app
