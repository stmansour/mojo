TOP=../..
BINDIR=${TOP}/tmp/mojo
THISDIR="msgbody"
COUNTOL=${TOP}/tools/bashtools/countol.sh

msgbody: *.go config.json
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	@${COUNTOL} golint
	go build
	@rm -f fail
	@echo "*** ${THISDIR}: completed ${THISDIR} ***"

config.json:
	/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	cp confdev.json config.json

clean:
	go clean
	rm -f msgbody, ver.go, tmp *.log *.json
	@echo "*** ${THISDIR}: completed clean ***"

package: msgbody
	@echo "*** ${THISDIR}: completed package ***"


test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean msgbody package
