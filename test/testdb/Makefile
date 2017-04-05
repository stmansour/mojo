TOP=../..
BINDIR=../../tmp/mojo

testdb:
	@echo "**** build complete in testdb ***"

test:
	@touch fail
	./functest.sh
	mysqldump --no-defaults rentroll >restore.sql
	mv restore.sql ../ws/
	@rm -f fail

package:
	@echo "*** package completed in testdb ***"


clean:
	rm -f log mojolog llog ll.g err.txt [a-z] [a-z][a-z1-9] [a-z].g gold/*.g fail
	@echo "*** CLEAN completed in testdb ***"

