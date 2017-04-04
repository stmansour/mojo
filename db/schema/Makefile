TOP=../..

newdb:
	mysql --no-defaults < schema.sql
	mysql --no-defaults < testdb.sql

clean:
	@echo "*** CLEAN completed in schema ***"

test:
	@echo "*** TEST completed in schema ***"

install:
	cp schema.sql ../bin/

testdb:
	mysql --no-defaults < schema.sql
	mysql --no-defaults < testdb.sql
