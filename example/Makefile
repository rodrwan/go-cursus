VERSION=v0.0.1
SVC=broker
BIN=$(PWD)/bin/$(SVC)

GO ?= go
LDFLAGS='-extldflags "static" -X main.svcVersion=$(VERSION) -X main.svcName=$(SVC)'
TAGS=netgo -installsuffix netgo

.PHONY: users
users u:
	@echo "[running] Running users service..."
	@go run users/cmd/server/main.go

.PHONY: address
address a:
	@echo "[running] Running address service..."
	@go run address/cmd/server/main.go



