TOP=../..
COUNTOL=${TOP}/tools/bashtools/countol.sh

# This makefile creates a new db each time make is called
# for the build because if there are syntax errors in
# schema.sql, they tend not to be caught for a long time.
# Creating the empty db with the schema guarantees that
# the schema.sql file stays in good shape

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
