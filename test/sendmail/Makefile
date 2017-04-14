TOP=../..
BINDIR=${TOP}/tmp/mojo
THISDIR="sendmail"

sendmail: *.go
	go vet
	golint
	go build
	@echo "*** ${THISDIR}: completed ${THISDIR} ***"

clean:
	go clean
	rm -f sendmail, ver.go, tmp *.log
	@echo "*** ${THISDIR}: completed clean ***"

package: sendmail
	cp sendmail ${BINDIR}/mailtester
	cp testmsg.html ${BINDIR}/
	@echo "*** ${THISDIR}: completed package ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean sendmail package
