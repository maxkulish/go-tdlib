TAG := v1.5.0
ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_DIR = $(ROOT_DIR)/bin
MAC_BIN = ${BIN_DIR}/macOS/


help: _help_

_help_:
	@echo make schema-update - download td_api.tl file from the Github repository
	@echo make generate-json - convert td_api.tl to td_api.json
	@echo make generate-code - read tl file and generate Golang code


schema-update:
	curl https://raw.githubusercontent.com/tdlib/td/${TAG}/td/generate/scheme/td_api.tl 2>/dev/null > ./data/td_api.tl

generate-json:
	go run ./cmd/generate-json.go \
		-version "${TAG}" \
		-output "./data/td_api.json"

generate-code:
	go run ./cmd/generate-code.go \
		-version "${TAG}" \
		-outputDir "./client" \
		-package client \
		-functionFile function.go \
		-typeFile type.go \
		-unmarshalerFile unmarshaler.go
	go fmt ./...


build_mac:
	cd $(ROOT_DIR)
	mkdir -p ${MAC_BIN}
	GOOS=darwin GOARCH=amd64 go build -o ${MAC_BIN}/telega ./main.go