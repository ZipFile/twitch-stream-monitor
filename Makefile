.DEFAULT_GOAL := build

build:
	go build

fmt:
	go fmt ./...

imports:
	find internal -name '*.go' -exec goimports -w {} \;
	goimports -w *.go

tests:
	go test ./... -coverprofile=coverage.out

tests-short:
	go test -short ./... -coverprofile=coverage.out

coverage:
	go tool cover -html=coverage.out
