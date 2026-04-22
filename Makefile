# PicPrism Makefile

GO=go
GOPATH_BIN=$(shell go env GOPATH)/bin

.PHONY: build dev tidy

build:
	cd web && npm run build
	$(GO) build -o picprism ./cmd/picprism

tidy:
	$(GO) mod tidy

dev-backend:
	PICPRISM_DATA_DIR=./data PICPRISM_TOKEN=dev_token $(GO) run ./cmd/picprism

dev-frontend:
	cd web && npm run dev

docker-build:
	docker build -t picprism:latest .
