TOP=../..
THIS=charnorm
BINDIR=${TOP}/tmp/mojo
COUNTOL=${TOP}/tools/bashtools/countol.sh

charnorm: *.go
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	@${COUNTOL} golint
	go build
	@rm -f fail
	@echo "*** ${THIS}: completed ***"

clean:
	go clean
	rm -f ${THIS}
	@echo "*** ${THIS}: completed clean ***"

package:
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	cp ${THIS} ${TOP}/tmp/mojo/${THIS}
	@rm -f fail
	@echo "*** ${THIS}: completed package ***"

test:
	@echo "*** ${THIS}: completed test ***"


all: clean csv package test
secure:
	@rm -f config.json confdev.json confprod.json
