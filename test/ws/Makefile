all:
	@echo "*** Completed in ws ***"

clean:
	rm -f mlog rawcmdout *.log log llog *.g ./gold/*.g err.txt [a-z] [a-z][a-z0-9] [a-z][0-9][0-9] qq? fail request smalldb.sql
	@echo "*** CLEAN completed in ws ***"

test:
	@touch fail
	#if [ ! -f smalldb.sql ]; then pushd ../testdb;make smalldb;cp smalldb.sql ../ws/;popd; fi
	pushd ../testdb;make smalldb;cp smalldb.sql ../ws/;popd
	./functest.sh
	@echo "*** TEST completed in ws ***"
	@rm -f fail

smalldb: testdb

testdb:
	cd ../testdb;make smalldb

package:
	@echo "*** PACKAGE completed in ws ***"

secure:
	@rm -f config.json confdev.json confprod.json
