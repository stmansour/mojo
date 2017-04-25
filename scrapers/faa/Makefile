TOP=../..
THISDIR="faa"
BINDIR=${TOP}/tmp/mojo
COUNTOL=${TOP}/tools/bashtools/countol.sh

faa: *.go
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	@${COUNTOL} golint
	go build
	@rm -f fail
	@echo "*** ${THISDIR}: completed faa ***"

clean:
	go clean
	rm -f faa
	@echo "*** ${THISDIR}: completed clean ***"

package:
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	cp faa ${TOP}/tmp/mojo/scrapefaa
	@rm -f fail
	@echo "*** ${THISDIR}: completed package ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

q:
	@echo "Create a new db in quick mode..."
	@${BINDIR}/mojonewdb
	@./faa -q
	@mysqldump --no-defaults mojo >faa.sql

all: clean faa package test
