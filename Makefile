SHELL=/usr/bin/env bash

.PHONY: clean
clean:
	rm claim-punk

.PHONY: all
all:
	go build -o claim-punk *.go