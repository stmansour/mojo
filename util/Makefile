TOP=..
THISDIR="util"
COUNTOL=${TOP}/tools/bashtools/countol.sh

util: *.go
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	@${COUNTOL} golint
	go build
	go install
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
