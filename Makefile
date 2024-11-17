dev:
	go run cmd/main.go

build:
	go build -ldflags "-X main.version=1.1.0" -o go_vercel_cli cmd/main.go

install_local:
	@mkdir -p /usr/local/bin/go_vercel_cli # Ensure the directory exists
	@mv go_vercel_cli /usr/local/bin/go_vercel_cli/go_vercel_cli

install: build install_local
	@echo "Build complete"
	@echo "Installing locally"

.PHONY: dev build install_local install
