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
	go generate ./...

build: generate
	pnpm run build
	go build -o app ./cmd/app
