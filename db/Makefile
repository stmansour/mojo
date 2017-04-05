TOP=..
DIRS=schema

db: *.go
	for dir in $(DIRS); do make -C $$dir;done
	go vet
	golint
	go build
	go install

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	go clean

package:
	go install
	for dir in $(DIRS); do make -C $$dir package;done

all: clean db package
