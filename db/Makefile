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

install:
	if [ -d ./bin ]; then rm -rf ./bin; fi
	mkdir ./bin
	for dir in $(DIRS); do make -C $$dir install;done

all: clean db
