TOP=..
WEBDOCDIR=${TOP}/../webdoc
WEBDOC=${WEBDOCDIR}/webdoc
WEBCSS=${WEBDOCDIR}/doc.css
RELDIR=${TOP}/tmp/mojo
THISDIR="mojosrv"

mojosrv: *.go ver.go webdoc
	go vet
	golint
	go build
	@echo "*** ${THISDIR}: completed mojosrv ***"

webdoc:
	rm -rf doc
	mkdir doc
	${WEBDOC}
	@cp ${WEBDOCDIR}/doc.css doc/
	@cp ${WEBDOCDIR}/tester.html doc/

ver.go:
	./mkver.sh

clean:
	go clean
	rm -f mojosrv, ver.go, tmp *.log
	@echo "*** ${THISDIR}: completed clean ***"

package:
	mkdir -p ${RELDIR}
	cp mojosrv ${RELDIR}/mojosrv
	@echo "*** ${THISDIR}: completed package ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean mojosrv package
