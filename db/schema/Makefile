TOP=../..

newdb:
	mysql --no-defaults < schema.sql

clean:
	@echo "*** CLEAN completed in schema ***"

test:
	@echo "*** TEST completed in schema ***"

package:
	cp schema.sql ${TOP}/tmp/mojo/

testdb:
	mysql --no-defaults < schema.sql
	mysql --no-defaults < testdb.sql
