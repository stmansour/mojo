TOP=../..
THISDIR="mojonewdb"
TOOL="mojonewdb"
COUNTOL=${TOP}/tools/bashtools/countol.sh

mojonewdb: *.sh

clean:
	rm -f ${THISDIR} fail conf*.json
	@echo "*** Clean completed in admin/${THISDIR} ***"

install: mojonewdb
	@echo "*** install completed in admin/${THISDIR} ***"

package: mojonewdb
	touch fail
	cp ${TOOL} ${TOP}/tmp/mojo/
	@echo "*** Packaging completed in admin/${THISDIR} ***"
	@rm -f fail

test:
	@echo "*** Testing completed in admin/${THISDIR} ***"
