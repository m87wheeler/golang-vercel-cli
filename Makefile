.PHONY: dev
dev:
	go run cmd/main.go

.PHONY: build
build:
	go build cmd/main.go 