TOP=..
BINDIR=${TOP}/tmp/mojo
THISDIR="mailsend"
COUNTOL=${TOP}/tools/bashtools/countol.sh

mailsend: *.go
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build
	go install
	@rm -f fail
	@echo "*** ${THISDIR}: completed ${THISDIR} ***"

clean:
	go clean
	rm -f mailsend, ver.go, tmp *.log conf*.json
	@echo "*** ${THISDIR}: completed clean ***"

package:
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	#insert package commands here
	@rm -f fail
	@echo "*** ${THISDIR}: completed package ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean mailsend package
secure:
	@rm -f config.json confdev.json confprod.json
