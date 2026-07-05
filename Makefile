BINARY=sysadmin-cli

build:
	go build -o $(BINARY) .

install: build
	cp $(BINARY) /usr/local/bin/

clean:
	rm -f $(BINARY)

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

.PHONY: build install clean test vet lint tidy
