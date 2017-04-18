TOP=..
BINDIR=${TOP}/tmp/mojo
THISDIR="mailsend"

mailsend: *.go
	go vet
	golint
	go build
	go install
	@echo "*** ${THISDIR}: completed ${THISDIR} ***"

clean:
	go clean
	rm -f mailsend, ver.go, tmp *.log
	@echo "*** ${THISDIR}: completed clean ***"

package:
	@echo "*** ${THISDIR}: completed package ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean mailsend package
