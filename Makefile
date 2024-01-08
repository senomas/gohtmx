.PHONY: FORCE
SHELL=/bin/bash
export

test: FORCE
	go test -v ./...
