TOP=..
DIRS=schema
THISDIR=db
COUNTOL=${TOP}/tools/bashtools/countol.sh

db: *.go
	for dir in $(DIRS); do make -C $$dir;done
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build
	go install
	@rm -f fail
	@echo "*** ${THISDIR}: completed db ***"

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	go clean
	@echo "*** ${THISDIR}: completed clean ***"

test:
	@echo "*** ${THISDIR}: completed test ***"


package:
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	go install
	for dir in $(DIRS); do make -C $$dir package;done
	@rm -f fail
	@echo "*** ${THISDIR}: completed package ***"

all: clean db package
secure:
	for dir in $(DIRS); do make -C $${dir} secure;done
	@rm -f config.json confdev.json confprod.json
