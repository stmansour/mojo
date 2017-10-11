#!/bin/bash

#---------------------------------------------------------------
# TOP is the directory where RentRoll begins. It is used
# in base.sh to set other useful directories such as ${BASHDIR}
#---------------------------------------------------------------
TOP=../..
BINDIR=${TOP}/tmp/mojo


TESTNAME="Web Services"
TESTSUMMARY="Test Web Services"

CREATENEWDB=0

#---------------------------------------------------------------
#  Use the testdb for these tests...
#---------------------------------------------------------------
echo "Create new database..." 
# mysql --no-defaults mojo < smalldb.sql
pushd ../testdb
echo "Loading bigdb"
make bigdb
popd

source ../share/base.sh
GRP="FAA Tech Ops"
echo "${BINDIR}/mojocsv -g \"${GRP}\" -cg -f faatechops.csv"
${BINDIR}/mojocsv -g "${GRP}" -o -cg -f faatechops.csv  >log 2>&1

logcheck
