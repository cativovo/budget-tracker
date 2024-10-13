seed:
	go run ./cmd/seed

cleandb:
	go run ./cmd/seed -c

compile:
	make -j3 compiletailwind compiletempl compileesbuild

compiletailwind:
	pnpm run tailwind:compile

compileesbuild:
	pnpm run esbuild:compile

compiletempl:
	templ generate

liveserver:
	go run github.com/air-verse/air@v1.60.0 \
	--build.cmd "go build -o ./tmp/bin/main ./cmd/app" \
	--build.bin "./tmp/bin/main" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

livecompile:
	go run github.com/air-verse/air@v1.60.0 \
	--build.cmd "make compile" \
	--build.bin "true" \
	--build.exclude_dir "assets" \
	--build.include_ext "js,css,templ"

live:
	make -j2 livecompile liveserver

generate:
	go generate ./...
