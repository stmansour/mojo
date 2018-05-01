TOP=../..
BINDIR=../../tmp/mojo
COUNTOL=${TOP}/tools/bashtools/countol.sh

testdb:
	@echo "**** build complete in testdb ***"

# bigdb is the full faa database
bigdb:
	@gunzip bigdb.sql.gz
	mysql --no-defaults mojo <bigdb.sql
	@gzip bigdb.sql

# smalldb is made up of the portion of the FAA directory
# where last names begin with 'aa', plus the information
# added to the database by the mailsend client (run with
# -n to just update the db andavoid actually sending mail).
smalldb:
	if [ ! -f faa.sql ]; then cp ${TOP}/scrapers/faa/faa.sql .; fi
	mysql --no-defaults mojo <faa.sql
	cd ${TOP}/test/mailsend;./mailsend -setup
	@mysqldump --no-defaults mojo >smalldb.sql

# ws (web service) testing is done with the smalldb
test:
	@touch fail
	mysqldump --no-defaults mojo  >restore.sql
	mv restore.sql ../ws/
	@rm -f fail

package:
	@echo "*** package completed in testdb ***"


clean:
	rm -f log mojolog llog ll.g err.txt [a-z] [a-z][a-z1-9] [a-z].g gold/*.g fail *.sql
	@echo "*** CLEAN completed in testdb ***"

snapshot:
	@echo "*** snapshot current database as smalldb"
	@mysqldump --no-defaults mojo >smalldb.sql
	@gzip -f smalldb.sql
secure:
	@rm -f config.json confdev.json confprod.json
