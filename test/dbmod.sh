#!/bin/bash

#==========================================================================
#  This script performs SQL schema changes on the test databases that are
#  saved as SQL files in the test directory. It loads them, performs the
#  ALTER commands, then saves the sql file.
# 
#  If the test file uses its own database saved as a .sql file, make sure
#  it is listed in the dbs array
#==========================================================================

MODFILE="dbqqqmods.sql"
MYSQL="mysql --no-defaults"
MYSQLDUMP="mysqldump --no-defaults"
DBNAME="mojo"

#=====================================================
#  Put modifications to schema in the lines below
#=====================================================
cat >${MODFILE} <<EOF
ALTER TABLE People ADD Email2 varchar(5) NOT NULL DEFAULT '' AFTER Email1;
EOF

#=====================================================
#  Put dir/sqlfilename in the list below
#=====================================================
declare -a dbs=(
	'../scrapers/faa/faa.sql'
	'testdb/bigdb.sql'
	'csv2/smalldb.sql'
	'ws/restore.sql'
)

pushd testdb
gunzip bigdb.sql.gz
popd

for f in "${dbs[@]}"
do
	echo -n "${f}: loading... "
	${MYSQL} ${DBNAME} < ${f}
	echo -n "updating... "
	${MYSQL} ${DBNAME} < ${MODFILE}
	echo -n "saving... "
	${MYSQLDUMP} ${DBNAME} > ${f}
	echo "done"
done

pushd testdb
gzip bigdb.sql
popd
