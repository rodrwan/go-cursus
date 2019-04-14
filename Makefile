VERSION=v0.0.1
SVC=broker
BIN=$(PWD)/bin/$(SVC)

GO ?= go
LDFLAGS='-extldflags "static" -X main.svcVersion=$(VERSION) -X main.svcName=$(SVC)'
TAGS=netgo -installsuffix netgo

run r:
	@echo "[running] Running service..."
	@go run cmd/broker/*.go

build b:
	@echo "[build] Building service..."
	@cd cmd/broker && $(GO) build -o $(BIN) -ldflags=$(LDFLAGS) -tags $(TAGS)

linux l:
	@echo "[build] Building for linux..."
	@cd cmd/broker && GOOS=linux $(GO) build -a -o $(BIN) --ldflags $(LDFLAGS) -tags $(TAGS)
