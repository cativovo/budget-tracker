seed:
	go run ./cmd/seed

cleandb:
	go run ./cmd/seed -c

livetempl:
	templ generate --watch

liveserver:
	go run github.com/air-verse/air@v1.60.0 \
	--build.cmd "go build -o ./tmp/bin/main ./cmd/app" --build.bin "./tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

live:
	make -j2 livetempl liveserver

sqlcgenerate:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate
