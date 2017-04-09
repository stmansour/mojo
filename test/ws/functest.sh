#!/bin/bash

#---------------------------------------------------------------
# TOP is the directory where RentRoll begins. It is used
# in base.sh to set other useful directories such as ${BASHDIR}
#---------------------------------------------------------------
TOP=../..

TESTNAME="Web Services"
TESTSUMMARY="Test Web Services"

CREATENEWDB=0

#---------------------------------------------------------------
#  Use the testdb for these tests...
#---------------------------------------------------------------
# echo "Create new database..." 
# mysql --no-defaults rentroll < restore.sql

source ../share/base.sh

echo "STARTING MOJO SERVER"
startMojoServer

curl http://localhost:8275/v1/ping/

# Get Specificy PaymentType 
echo "request%3D%7B%22cmd%22%3A%22get%22%2C%22selected%22%3A%5B%5D%2C%22limit%22%3A100%2C%22offset%22%3A0%7D" > request
dojsonPOST "http://localhost:8275/v1/people/1" "request" "y"  "WebService--GetPeople"
dojsonGET "http://localhost:8275/v1/peoplecount/" "a"  "WebService--PeopleCount"
dojsonGET "http://localhost:8275/v1/groupcount/" "b"  "WebService--GroupCount"
dojsonGET "http://localhost:8275/v1/peoplestats/" "c"  "WebService--PeopleStats"

echo "request%3D%7B%22cmd%22%3A%22%22%2C%22limit%22%3A100%2C%22offset%22%3A0%7D" > request
dojsonPOST "http://localhost:8275/v1/queries/" "request" "d"  "WebService--SearchQueries"

dojsonGET "http://localhost:8275/v1/qrescount/1" "e"  "WebService--QueryResultsCount"


stopMojoServer
echo "RENTROLL SERVER STOPPED" 

logcheck
