TOP=..
util: *.go
	go vet
	golint
	go build
	go install

clean:
	go clean

install:

all: clean util
