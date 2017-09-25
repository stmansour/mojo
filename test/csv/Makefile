TOP=../..
BINDIR=${TOP}/tmp/mojo
THISDIR=csv
all:
	@echo "*** Completed in ${THISDIR} ***"

clean:
	rm -f *.log log llog *.g ./gold/*.g err.txt [a-z] [a-z][a-z1-9] qq? fail request smalldb.sql rawcmdout
	@echo "*** CLEAN completed in ${THISDIR} ***"

test:
	@touch fail
	#if [ ! -f smalldb.sql ]; then pushd ../testdb;make smalldb;cp smalldb.sql ../${THISDIR}/;popd; fi
	pushd ../testdb;make smalldb;cp smalldb.sql ../${THISDIR}/;popd
	./functest.sh
	@echo "*** TEST completed in ${THISDIR} ***"
	@rm -f fail

smalldb: testdb

testdb:
	cd ../testdb;make smalldb

package:
	@echo "*** PACKAGE completed in ${THISDIR} ***"

