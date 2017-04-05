TOP=../..

faa: *.go
	go vet
	golint
	go build

clean:
	go clean

package:
	cp faa ${TOP}/tmp/mojo/scrapefaa

all: clean faa
