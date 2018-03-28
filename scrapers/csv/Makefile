TOP=../..
THISDIR="csv"
BINDIR=${TOP}/tmp/mojo
COUNTOL=${TOP}/tools/bashtools/countol.sh

csv: *.go config.json
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	@${COUNTOL} golint
	go build
	@rm -f fail
	@echo "*** ${THISDIR}: completed ***"

config.json:
	/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	cp confdev.json config.json

clean:
	go clean
	rm -f csv conf*.json *.log g.csv g[0-9].csv
	@echo "*** ${THISDIR}: completed clean ***"

package:
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	cp csv ${TOP}/tmp/mojo/mojocsv
	@rm -f fail
	@echo "*** ${THISDIR}: completed package ***"

test:
	@echo "*** ${THISDIR}: completed test ***"


all: clean csv package test
