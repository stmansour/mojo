TOP=..
WEBDOCDIR=${TOP}/../webdoc
WEBDOC=${WEBDOCDIR}/webdoc
WEBCSS=${WEBDOCDIR}/doc.css
BINDIR=${TOP}/tmp/mojo
THISDIR="mojosrv"

mojosrv: *.go ver.go webdoc conf.json
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

conf.json:
	cp confdev.json conf.json

ver.go:
	./mkver.sh

clean:
	go clean
	rm -f mojosrv, ver.go, tmp *.log
	@echo "*** ${THISDIR}: completed clean ***"

package:
	mkdir -p ${BINDIR}
	cp mojosrv ${BINDIR}/mojosrv
	cp conf.json ${BINDIR}/
	cp -r doc ${BINDIR}/
	cp -r html ${BINDIR}/
	@echo "*** ${THISDIR}: completed package ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean mojosrv package

images:
	tar cvf images.tar html/images
	gzip mojoimages.tar
