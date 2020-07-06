#!/bin/bash

#---------------------------------------------------------------
# TOP is the directory where RentRoll begins. It is used
# in base.sh to set other useful directories such as ${BASHDIR}
#---------------------------------------------------------------
BINDIR=../../tmp/mojo

TESTNAME="Imports"
TESTSUMMARY="Test Import Functions"
NEWDB="${BINDIR}/mojonewdb"
GRP="smanmusic"

CREATENEWDB=1

source ../share/base.sh

showResults() {
	echo "Test ${TFILES}:  ${PF}"
	if [ "${PF}" = "pass" ]; then
		passmsg
	else
		failmsg
	fi
}

#---------------------------------------------------------------
#  Use the testdb for these tests...
#---------------------------------------------------------------
echo "Create new database..."
${NEWDB}
if [ ${CREATENEWDB} != "0" ]; then
	mysql --no-defaults mojo < xa.sql
fi

#------------------------------------------------------------------------------
#  We start with xa.sql.  It has 2 entries. In this test module we test many
#  ways to update these entries.
#
#  The CSV library looks for the following fields in the second line of the
#  CSV file:
#     0  FirstName
#     1  MiddleName
#     2  LastName
#     3  PreferredName
#     4  JobTitle
#     5  OfficePhone
#     6  OfficeFax
#     7  Email1
#     8  Email2
#     9  MailAddress
#    10  MailAddress2
#    11  MailCity
#    12  MailState
#    13  MailPostalCode
#    14  MailCountry
#    15  RoomNumber
#    16  MailStop
#
#  Map any fields you're trying to import accordingly
#------------------------------------------------------------------------------


#------------------------------------------------------------------------------
#  TEST a
#  Verify that an existing entry is updatedwith new data that comes from the
#  CSV file.
#
#  Scenario:
#  xa.csv contains updated MiddleName information. The MiddleName in the
#  existing entry is empty. Verify that the MiddleName info is correctly
#  copied into the database.
#
#  Expected Results:
#   1.  MiddleName should be added to PID 1: "BillyBob". This should result
#       in Updated Entries = 1.
#------------------------------------------------------------------------------
TFILES="a"
echo "SINGLETEST, TFILES = ${SINGLETEST}${TFILES}"
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
	TESTNAME="Update Missing Info"
	TESTCOUNT=1
	${BINDIR}/mojocsv -g "${GRP}" -cg -f xa.csv > "x${TFILES}"
	RES=$(grep "Updated Entries:" xa | sed 's/Updated Entries: *//')
	PF="fail"
	if [ "${RES}" = "1" ]; then
		PF="pass"
	fi
	echo "Test ${TFILES}:  ${PF}"
	showResults
fi

#------------------------------------------------------------------------------
#  TEST b
#  Verify that updating an existing entry with exactly the same information
#  that exists in the database already basically does not make any changes to
#  the db
#
#  Scenario:
#  re-import the same csv file as in test a, this should result in nothing
#  happening to the database.
#
#  Expected Results:
#   1.  No changes are made to the database
#------------------------------------------------------------------------------
TFILES="b"
echo "SINGLETEST, TFILES = ${SINGLETEST}${TFILES}"
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
	TESTNAME="Skip DB Write If No Change"
	TESTCOUNT=1
	${BINDIR}/mojocsv -g "${GRP}" -cg -f xa.csv > "x${TFILES}"
	RES=$(grep "Updated Entries:" x${TFILES} | sed 's/Updated Entries: *//')
	PF="fail"
	if [ "${RES}" = "0" ]; then
		PF="pass"
	fi
	echo "Test ${TFILES}:  ${PF}"
	showResults
fi

#------------------------------------------------------------------------------
#  TEST c
#  Verify that updating an existing entry with a filled in field where the
#  csv file has an empty field does not update the existing entry
#
#  Scenario:
#  re-import the a csv file that has all the same information for Shannon
#  except that the MiddleName is now empty in the csv file, this should result
#  no db updates.
#
#  Expected Results:
#   1.  No changes are made to the database
#------------------------------------------------------------------------------
TFILES="c"
echo "SINGLETEST, TFILES = ${SINGLETEST}${TFILES}"
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
	#-------------------------------------------------
	# first , guarantee that the middle name is added
	#-------------------------------------------------
	TESTNAME="Blank Info Does Not Overwrite"
	TESTCOUNT=1
	${BINDIR}/mojocsv -g "${GRP}" -cg -f xa.csv > "x${TFILES}"
	RES=$(grep "Updated Entries:" xa | sed 's/Updated Entries: *//')
	PF="fail"
	if [ "${RES}" != "1" ]; then
		PF="fail"
		echo "Test ${TFILES}: ${PF} -- did not add MiddleName for initial condition"
		exit 1
	fi
	#-----------------------------------------------------------
	# now read a csv file that has no MiddleName for Shannon
	#-----------------------------------------------------------
	((TESTCOUNT++))
	${BINDIR}/mojocsv -g "${GRP}" -cg -f xc.csv > "x${TFILES}"
	RES=$(grep "Updated Entries:" x${TFILES} | sed 's/Updated Entries: *//')
	if [ "${RES}" = "0" ]; then
		PF="pass"
	fi
	showResults
fi

#------------------------------------------------------------------------------
#  TEST d
#  Verify that updating an existing entry with a filled in field where the
#  csv file has an different value for the field gets the update to the db
#
#  Scenario:
#  Shannon will have a MiddleName of BillyBob to start with. Then when we
#  import xd.csv it should update the MiddleName to CornDog
#
#  Expected Results:
#   1.  When test finishes, Shannon should have MiddleName CornDog
#------------------------------------------------------------------------------
TFILES="d"
echo "SINGLETEST, TFILES = ${SINGLETEST}${TFILES}"
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
	#-------------------------------------------------
	# first , guarantee that the middle name is added
	#-------------------------------------------------
	TESTNAME="Updated Info Overwrites"
	TESTCOUNT=1
	${BINDIR}/mojocsv -g "${GRP}" -cg -f xa.csv > "x${TFILES}"
	RES=$(grep "Updated Entries:" xa | sed 's/Updated Entries: *//')
	PF="fail"
	if [ "${RES}" != "1" ]; then
		PF="fail"
		echo "Test ${TFILES}: ${PF} -- did not add MiddleName for initial condition"
		exit 1
	fi
	#-----------------------------------------------------------
	# now read a csv file that has no MiddleName for Shannon
	#-----------------------------------------------------------
	((TESTCOUNT++))
	${BINDIR}/mojocsv -g "${GRP}" -cg -f xd.csv > "x${TFILES}"
	RES=$(grep "Updated Entries:" x${TFILES} | sed 's/Updated Entries: *//')
	if [ "${RES}" = "1" ]; then
		PF="pass"
	fi
	echo "Test ${TFILES}:  ${PF}"
	showResults
fi
