TOP=../..
THISDIR="faa"
BINDIR=${TOP}/tmp/mojo

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

q:
	@echo "Create a new db in quick mode..."
	@${BINDIR}/mojonewdb
	@./faa -q
	@mysqldump --no-defaults mojo >faa.sql

all: clean faa package test
