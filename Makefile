dev:
	go run cmd/main.go

build:
	go build -ldflags "-X main.version=1.0.0" -o go_vercel_cli cmd/main.go

.PHONY: dev build