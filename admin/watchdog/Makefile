TOP=../..
THISDIR="watchdog"
TOOL="mojowatchdog"
COUNTOL=${TOP}/tools/bashtools/countol.sh

mojowatchdog: mojowatchdog.sh

clean:
	rm -f ${TOOL} fail conf*.json
	@echo "*** Clean completed in admin/${THISDIR} ***"

install: ${TOOL}
	@echo "*** install completed in admin/${THISDIR} ***"

package: mojowatchdog
	touch fail
	cp ${TOOL} ${TOP}/tmp/mojo/
	@echo "*** Packaging completed in admin/${THISDIR} ***"
	rm -f fail

test:
	@echo "*** Testing completed in admin/${THISDIR} ***"
