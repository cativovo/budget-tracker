ENV_FILE ?= .env.development.local
include $(ENV_FILE)
export

seed:
	go run ./cmd/seed

cleandb:
	go run ./cmd/seed -c

liveui:
	pnpm --dir ./ui run dev

liveserver:
	go run github.com/air-verse/air@v1.60.0 \
	--build.cmd "go build -o ./tmp/bin/main ./cmd/app" \
	--build.bin "./tmp/bin/main" \
	--build.include_dir "cmd/app,internal" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

live:
	make -j2 liveui liveserver

test:
	go test ./...

build:
	pnpm --dir ./ui run build
	go build -o ./bin/app ./cmd/app

start:
	BUDGET_TRACKER_ENV=production ./bin/app
