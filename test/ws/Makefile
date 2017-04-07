all:
	@echo "*** Completed in ws ***"

clean:
	rm -f *.log log llog *.g ./gold/*.g err.txt [a-z] [a-z][a-z1-9] qq? fail request
	@echo "*** CLEAN completed in ws ***"

test:
	@touch fail
	mysql --no-defaults mojo < restore.sql
	./functest.sh
	@echo "*** TEST completed in ws ***"
	@rm -f fail

package:
	@echo "*** PACKAGE completed in ws ***"

backupdb:
	mysqldump --no-defaults mojo > restore.sql

restoredb:
	mysql --no-defaults mojo < restore.sql
