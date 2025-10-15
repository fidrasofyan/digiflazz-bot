BINARY  := bin/digiflazz-bot
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -s -w \
	-X 'main.AppVersion=$(VERSION)' \

build:
	staticcheck ./...
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/digiflazz-bot
