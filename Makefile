SHELL=/usr/bin/env bash

GOCC?=go

.PHONY: clean
clean:
	rm claim-punk

.PHONY: all
all:
	go build -o claim-punk *.go

type-gen: 
	$(GOCC) run ./gen/main.go