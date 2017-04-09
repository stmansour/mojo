TOP=../..
BINDIR=../../tmp/mojo

testdb:
	@echo "**** build complete in testdb ***"

bigdb:
	@gunzip bigdb.sql.gz
	mysql --no-defaults mojo <bigdb.sql
	@gzip bigdb.sql

smalldb:
	@gunzip smalldb.sql.gz
	mysql --no-defaults mojo <smalldb.sql
	@gzip smalldb.sql

test:
	@touch fail
	mysqldump --no-defaults mojo  >restore.sql
	mv restore.sql ../ws/
	@rm -f fail

package:
	@echo "*** package completed in testdb ***"


clean:
	rm -f log mojolog llog ll.g err.txt [a-z] [a-z][a-z1-9] [a-z].g gold/*.g fail
	@echo "*** CLEAN completed in testdb ***"

snapshot:
	@echo "*** snapshot current database as smalldb"
	@mysqldump --no-defaults mojo >smalldb.sql
	@gzip -f smalldb.sql
