TOP=../..
BINDIR=${TOP}/tmp/mojo
THISDIR="mailsend"

mailsend: *.go
	go vet
	golint
	go build
	@echo "*** ${THISDIR}: completed ${THISDIR} ***"

clean:
	go clean
	rm -f mailsend, ver.go, tmp *.log
	@echo "*** ${THISDIR}: completed clean ***"

package: mailsend
	cp mailsend ${BINDIR}/mailtester
	cp testmsg.html mscript.sh resetdb.sh ${BINDIR}/
	@echo "*** ${THISDIR}: completed package ***"

newdb:
	${BINDIR}/mojonewdb
	./mailsend -setup

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean mailsend package
