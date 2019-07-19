#!/usr/bin/make -f

test:
	go test -timeout=1s -short ./...

compile:
	go build ./...

build: test compile

document:
	go install github.com/robertkrimen/godocdown/godocdown \
		&& godocdown > README.md

.PHONY: test compile build document
