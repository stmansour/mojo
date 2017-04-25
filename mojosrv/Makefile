TOP=..
WEBDOCDIR=${TOP}/../webdoc
WEBDOC=${WEBDOCDIR}/webdoc
WEBCSS=${WEBDOCDIR}/doc.css
BINDIR=${TOP}/tmp/mojo
THISDIR="mojosrv"
COUNTOL=${TOP}/tools/bashtools/countol.sh

mojosrv: ver.go *.go webdoc config.json
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	@${COUNTOL} golint
	go build
	@rm -f fail
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
	rm -f mojosrv ver.go tmp *.log *.json
	@echo "*** ${THISDIR}: completed clean ***"

package:
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	mkdir -p ${BINDIR}
	cp mojosrv ${BINDIR}/mojosrv
	cp config.json ${BINDIR}/
	cp -r doc ${BINDIR}/
	cp -r html ${BINDIR}/
	if [ -d "./js" ]; then cp -r js ${BINDIR}/; fi
	@rm -f fail
	@echo "*** ${THISDIR}: completed package ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean mojosrv package

images:
	tar cvf images.tar html/images
	gzip mojoimages.tar

mkprod:
	/usr/local/accord/bin/getfile.sh accord/db/confprod.json
	cp confprod.json config.json
