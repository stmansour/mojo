TOP=../..
THISDIR="faa"

faa: *.go
	go vet
	golint
	go build
	@echo "*** ${THISDIR}: completed faa ***"

clean:
	go clean
	rm -f faa
	@echo "*** ${THISDIR}: completed clean ***"

package:
	cp faa ${TOP}/tmp/mojo/scrapefaa
	@echo "*** ${THISDIR}: completed package ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean faa package test
