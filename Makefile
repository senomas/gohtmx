.PHONY: FORCE
SHELL=/bin/bash
export

test: FORCE build
	go test -v ./...

build: FORCE
	~/go/bin/templ generate
