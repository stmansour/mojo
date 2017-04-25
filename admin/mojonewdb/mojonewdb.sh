#!/bin/bash
MYSQLOPTS="--no-defaults"
DBBACKUP="MojoBackupDB"
DBNAME="mojo"
if [ -f /usr/local/bin/mysql ]; then
	MYSQL="/usr/local/bin/mysql"
elif [ -f /usr/bin/mysql ]; then
	MYSQL="/usr/bin/mysql"
else
	MYSQL="mysql"
fi

function MakeProdDB() {
	MYSQLOPTS=

	HOST=$(grep "MojoDbhost" config.json | sed -e 's/.*"MojoDbhost"[ \t]*:[ \t]*"//' | sed -e 's/",$//')
	PORT=$(grep "MojoDbport" config.json | sed -e 's/.*"MojoDbport"[ \t]*:[ \t]*//' | sed -e 's/[ \t]*,[ \t]*$//')
	MYSQLDUMP="${MYSQL}dump"
	${MYSQLDUMP} -h ${HOST} -P ${PORT} ${DBNAME} >${DBBACKUP}
	RC=$?
	if [ $RC == 0 ]; then
		echo "Mojo database exists. A backup was made to ${DBBACKUP}."
	else
		${MYSQL} -h ${HOST} -P ${PORT} <schema.sql
		echo "Mojo database does not exist. A new db was created."
	fi
}

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

pushd ${DIR}
PROD=$(grep '"Env"' config.json | grep 0 | wc -l)   # if this is production then PROD == 1, otherwise PROD == 0
if [ "${PROD}" = "1" ]; then
	MakeProdDB
else
	mysql ${MYSQLOPTS} <schema.sql
fi
popd

