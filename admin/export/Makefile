TOP=../..
BINDIR=${TOP}/tmp/mojo
THISDIR="export"
COUNTOL=${TOP}/tools/bashtools/countol.sh

export: *.go config.json
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build
	@rm -f fail
	@echo "*** ${THISDIR}: completed ${THISDIR} ***"

config.json:
	/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	cp confdev.json config.json

clean:
	go clean
	rm -f mailsend, ver.go, tmp *.log *.json
	@echo "*** ${THISDIR}: completed clean ***"

package: export
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	cp export ${BINDIR}/mojoexport
	@rm -f fail
	@echo "*** ${THISDIR}: completed package ***"

newdb:
	${BINDIR}/mojonewdb
	${BINDIR}/mailsend -setup

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean mailsend package
secure:
	@rm -f config.json confdev.json confprod.json
