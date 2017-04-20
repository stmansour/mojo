TOP=..
WEBDOCDIR=${TOP}/../webdoc
WEBDOC=${WEBDOCDIR}/webdoc
WEBCSS=${WEBDOCDIR}/doc.css
BINDIR=${TOP}/tmp/mojo
THISDIR="mojosrv"

mojosrv: *.go ver.go webdoc config.json
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

config.json:
	/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	cp confdev.json config.json

ver.go:
	./mkver.sh

clean:
	go clean
	rm -f mojosrv, ver.go, tmp *.log *.json
	@echo "*** ${THISDIR}: completed clean ***"

package:
	mkdir -p ${BINDIR}
	cp mojosrv ${BINDIR}/mojosrv
	cp config.json ${BINDIR}/
	cp -r doc ${BINDIR}/
	cp -r html ${BINDIR}/
	if [ -d "./js" ]; then cp -r js ${BINDIR}/; fi
	@echo "*** ${THISDIR}: completed package ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean mojosrv package

images:
	tar cvf images.tar html/images
	gzip mojoimages.tar
