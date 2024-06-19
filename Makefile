.PHONY: build
build:
	go vet
	go test .
	go build .

.PHONY: install
install:
	go install

.PHONY: update
update:
	go get -u
	go mod tidy
