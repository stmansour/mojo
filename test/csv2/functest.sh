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
mysql --no-defaults mojo < smalldb.sql

source ../share/base.sh
GRP="AccordTest"
echo "${BINDIR}/mojocsv -g \"${GRP}\" -cg  -f accordtest.csv"
${BINDIR}/mojocsv -g "${GRP}" -o -cg -f accordtest.csv  >log 2>&1

logcheck
