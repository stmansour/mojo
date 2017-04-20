TOP=../..
BINDIR=${TOP}/tmp/mojo
THISDIR="mailsend"

mailsend: *.go config.json
	go vet
	golint
	go build
	@echo "*** ${THISDIR}: completed ${THISDIR} ***"

config.json:
	/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	cp confdev.json config.json

clean:
	go clean
	rm -f mailsend, ver.go, tmp *.log *.json
	@echo "*** ${THISDIR}: completed clean ***"

package: mailsend
	cp mailsend ${BINDIR}/mailsend
	cp testmsg.html mscript.sh resetdb.sh ${BINDIR}/
	@echo "*** ${THISDIR}: completed package ***"

newdb:
	${BINDIR}/mojonewdb
	./mailsend -setup

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean mailsend package
