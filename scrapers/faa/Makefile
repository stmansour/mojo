TOP=../..
THISDIR="faa"
BINDIR=${TOP}/tmp/mojo
COUNTOL=${TOP}/tools/bashtools/countol.sh

faa: *.go config.json
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build
	@rm -f fail
	@echo "*** ${THISDIR}: completed faa ***"

config.json:
	/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	cp confdev.json config.json

clean:
	go clean
	rm -rf faa conf*.json *.log csvdump
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

update:
	@echo "Create a new updated database"
	@${BINDIR}/mojonewdb
	@./faa

all: clean faa package test
secure:
	@rm -f config.json confdev.json confprod.json
