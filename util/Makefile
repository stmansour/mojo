TOP=..
THISDIR="util"
COUNTOL=${TOP}/tools/bashtools/countol.sh

util: *.go
	go get -v golang.org/x/text/secure/precis
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	if [[ -f "${GOPATH}/bin/golint" ]]; then golint; fi
	go build
	go install
	go test
	@rm -f fail
	@echo "*** ${THISDIR}: completed util ***"

clean:
	go clean
	@echo "*** ${THISDIR}: completed clean ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

package:
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	#insert package commands here
	@rm -f fail
	@echo "*** util completed package ***"

all: clean util
secure:
	@rm -f config.json confdev.json confprod.json
