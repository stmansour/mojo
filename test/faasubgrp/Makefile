TOP=../..
BINDIR=${TOP}/tmp/mojo
THISDIR=faasubgrp
all:
	@echo "*** Completed in ${THISDIR} ***"

clean:
	rm -f *.log log llog *.g ./gold/*.g err.txt [a-z] [a-z][a-z1-9] qq? fail request smalldb.sql rawcmdout
	@echo "*** CLEAN completed in ${THISDIR} ***"

test:
	@touch fail
	./functest.sh
	@echo "*** TEST completed in ${THISDIR} ***"
	@rm -f fail

package:
	@echo "*** PACKAGE completed in ${THISDIR} ***"

