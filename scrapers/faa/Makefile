TOP=../..

faa: *.go
	go vet
	golint
	go build
	go install

clean:
	go clean
	rm -f faa

install:

all: clean faa
