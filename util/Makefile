TOP=..
util: *.go
	go vet
	golint
	go build
	go install

clean:
	go clean

package:
	@echo "*** util completed package **"

all: clean util
