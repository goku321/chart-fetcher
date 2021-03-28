PROJECTNAME := $(shell basename "$(PWD)")
all : clean fmt test build run

build:
	@echo " > Building chart fetcher..."
	@go build $(LDFLAGS) -o $(PROJECTNAME)

test:
	go test -count=1 ./... -v

clean:
	go clean

fmt:
	go fmt ./...