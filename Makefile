# SPDX-FileCopyrightText: 2025 Florian Wilhelm
#
# SPDX-License-Identifier: MIT

all: format build test

format:
	gofumpt -w $$(find . -name '*.go')

build:
	go build -v ./...

test:
	go test -v ./...

update:
	go get -u
	go mod tidy

demo:
	go run . -input sample.csv -flat

clean:
	rm *ods
	rm *fods