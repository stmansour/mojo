#############################################################################
# Mojo test infrastructure base
#
# Include this script in your test script to get the base testing capabilities.
#
# Your script can override the following values:
#
#############################################################################

TOOLSDIR=${TOP}/tools
BASHDIR=${TOOLSDIR}/bashtools
TOP=../..

#############################################################################
# Set default values
#############################################################################
SECONDS=0
ERRFILE="err.txt"
UNAME=$(uname)
LOGFILE="log"
MYSQLOPTS=""
MYSQL=$(which mysql)
TESTCOUNT=0			## this is an internal counter, your external script should not touch it
SHOWCOMMAND=0
SCRIPTPATH=$(pwd -P)

if [ "x${MANAGESERVER}" = "x" ]; then
	MANAGESERVER=1
fi
if [ "x${CREATENEWDB}" = "x" ]; then
	CREATENEWDB=1
fi
if [ "x${MOJOPORT}" = "x" ]; then
	MOJOPORT="8275"
fi
if [ "x${MOJOBIN}" = "x" ]; then
	MOJOBIN="../../tmp/mojo"
else
	echo "MOJOBIN was pre-set to:  \"${MOJOBIN}\""
fi
TREPORT="${TOP}/test/testreport.txt"

MOJO="${MOJOBIN}/mojosrv -A"
CSVLOAD="${MOJOBIN}/rrloadcsv"
GOLD="./gold"

SKIPCOMPARE=0
FORCEGOOD=0
TESTCOUNT=0
ASKBEFOREEXIT=0

#############################################################################
#  This code ensures that mysql does not touch production databases.
#  The way identity is kept, default usage of mysql or mysqldump often
#  goes straight to the production databases.
#############################################################################
if [ "${UNAME}" == "Darwin" -o "${IAMJENKINS}" == "jenkins" ]; then
	MYSQLOPTS="--no-defaults"
fi

#############################################################################
# pause()
#   Description:
#		Ask the user how to proceed.
#
#   Params:
#       none
#############################################################################
pause() {
	echo
	echo
	read -p "Press [Enter] to continue, M to move ${2} to gold/${2}.gold, Q or X to quit..." x
	x=$(echo "${x}" | tr "[:upper:]" "[:lower:]")
	if [ "${x}" == "q" -o "${x}" == "x" ]; then
		if [ ${MANAGESERVER} -eq 1 ]; then
			echo "STOPPING MOJO SERVER"
			pkill mojosrv
		fi
		exit 0
	elif [[ ${x} == "m" ]]; then
		echo "********************************************"
		echo "********************************************"
		echo "********************************************"
		echo "cp ${1} gold/${1}.gold"
		cp ${1} gold/${1}.gold
	fi

}

app() {
	echo "command is:  ${MOJO}  ${1}"
	${MOJO} ${1}
}

usage() {
	cat <<EOF

SYNOPSIS
	$0 [-c -f -o -r]

	Mojo test script. Compare the output of each step to its associated
	.gold known-good output. If they miscompare, fail and stop the script.
	If they match, keep going until all tasks are completed.

OPTIONS
	-c  Show each command that was executed.

	-f  Executes all the steps of the test but does not compare the output
	    to the known-good files. This is useful when making a slight change
	    to something just to see how it will work.

	-m  Do not run any server mgmt commands. Typically, this is used to
		run the test commands against an already-running server.

	-n  Do not create a new database, use the current database and simply
	    add to it.

	-o  Regenerate the .gold files based on the output from this run. Only
	    use this option if you're sure the output is correct. This option
	    can be a huge time saver, but use it with caution. All .gold files
	    are maintained in the ./${GOLD}/ directory.

	-t  testname
		Run only the single test named testname.
		The tests are conventionally named as a single lower case letter
		(a-z), though this is not a requirement.
EOF
}

##########################################################################
# elapsedtime()
# Shows the number of seconds that was needed to run this script
##########################################################################
elapsedtime() {
	duration=$SECONDS
	msg="ElapsedTime: $(($duration / 60)) min $(($duration % 60)) sec"
	echo "${msg}" >>${LOGFILE}
	echo "${msg}"

}

passmsg() {
	printf "PASSED  %-20.20s  %-40.40s  %6d  \n" "${TESTDIR}" "${TESTNAME}" ${TESTCOUNT} >> ${TREPORT}
}

failmsg() {
	printf "FAILED  %-20.20s  %-40.40s  %6d  \n" "${TESTDIR}" "${TESTNAME}" ${TESTCOUNT} >> ${TREPORT}
}

forcemsg() {
	printf "FORCED  %-20.20s  %-40.40s  %6d  \n" "${TESTDIR}" "${TESTNAME}" ${TESTCOUNT} >> ${TREPORT}
}

tdir() {
	local IFS=/
	local p n m
	p=( ${SCRIPTPATH} )
	n=${#p[@]}
	m=$(( n-1 ))
	TESTDIR=${p[$m]}
}



#------------------------------------------------------------------------------
#  encodeURI encodes data so that it can be passed in a URI.  It
#      does essentially what Javascript's encodeURI does.
#
#  INPUTS
#  $1  The string to encode
#
#  RETURNS
#      the return value is the encoded string.
#
#  USAGE:
#      data=$(encodeURI "4%2F1%2F2019")
#------------------------------------------------------------------------------
encodeURI() {
  local string="${1}"
  local strlen=${#string}
  local encoded=""
  local pos c o

  for (( pos=0 ; pos<strlen ; pos++ )); do
     c=${string:$pos:1}
     case "$c" in
        [-_.~a-zA-Z0-9] ) o="${c}" ;;
        * )               printf -v o '%%%02x' "'$c"
     esac
     encoded+="${o}"
  done
  echo "${encoded}"
}

#------------------------------------------------------------------------------
#  encodeRequest is just like encodeURI except that it saves the output
#      into a file named "request"
#
#  INPUTS
#  $1  The string to encode
#
#  RETURNS
#      nothing, but the encoded string will be in a file named "request"
#------------------------------------------------------------------------------
encodeRequest() {
  local string="${1}"
  local strlen=${#string}
  local encoded=""
  local pos c o

  for (( pos=0 ; pos<strlen ; pos++ )); do
     c=${string:$pos:1}
     case "$c" in
        [-_.~a-zA-Z0-9] ) o="${c}" ;;
        * )               printf -v o '%%%02x' "'$c"
     esac
     encoded+="${o}"
  done
  echo "${encoded}" > request
}

#############################################################################
# incStep()
#   Description:
#		Increment the STEP variable.  It is encapsulated here because
#       there may be additional steps to perform in the future.
#
#   Params:
#       none
#############################################################################
incStep() {
	((STEP++))
}

#############################################################################
# domojotest()
#    The purpose of this routine is to call rrloadcsv with the
#     parameters supplied in $2 and send its output to a file
#     named $1. After trrloadcsv completes, the output in $1 will
#     be compared with the output in gold/$1.gold.  If there are
#     no diffs, then the test passes.  If there are diffs, then
#     it terminates execution of the script after doing
#     the following:
#
#        (a) Displays the diffs
#        (b) Displays the mv command to use if the newly generated
#            output is now correct and the gold/$1.gold file needs
#            to be updated.  You can just copy the command and paste
#            it into your command line.  Very handy
#        (c) Displays the full command it used to generate the output
#            in $1. This is very handy for reproducing a problem.
#
#     Additionally, there are some Environment Variables that cause
#     it to perform several functions that are very handy:
#
#        SKIPCOMPARE - ${SKIPCOMPARE} defaults to 0. As long as its
#            value is 0 the output in $1 will be compared to
#            gold/$1.gold .  However, there may be times where
#            you want the script to run to completion even if the
#            output miscompares with what is in gold/*  By convention,
#            all of my "functest.sh" scripts use the -f option to
#            set this value.
#
#        FORCEGOOD - ${FORCEGOOD} is set to 0 by default. If it is set
#            set to 1 it means that the output generated and stored in
#            $1 during this run is known to be "correct", even though
#            it may be different than what is in gold/$1.gold.  It will
#            automatically copy $1 to gold/$1.go. This is
#            extremely handy if a change was made to the table output
#            generator, or if any new fields were added to the database
#            and you've validated in some other way that everything is
#            working after such a change.  By convention, all of my
#            "function.sh" scripts use the -o option to set FORCEGOOD
#            to 1.
#
#	Parameters:
# 		$1 = base file name
#		$2 = app options to reproduce
# 		$3 = title for reporting.  No spaces
#############################################################################
########################################
# dorrtest()
#	Parameters:
# 		$1 = base file name
#		$2 = app options to reproduce
# 		$3 = title
########################################
dorrtest () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $1 $3
	${MOJO} ${2} >${1} 2>&1

	if [ "${FORCEGOOD}" = "1" ]; then
		cp ${1} ${GOLD}/${1}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${1}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${1}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi
		UDIFFS=$(diff ${1} ${GOLD}/${1}.gold | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${CSVLOAD} ${2}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED..." >> ${ERRFILE}
			echo "Differences in ${1} are as follows:" >> ${ERRFILE}
			diff ${GOLD}/${1}.gold ${1} >> ${ERRFILE}
			echo "If correct:  mv ${1} ${GOLD}/${1}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${MOJO} ${2}" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${3}
			else
				if [ ${MANAGESERVER} -eq 1 ]; then
					echo "STOPPING MOJO SERVER"
					pkill mojosrv
				fi
				exit 1
			fi
		fi
	else
		echo
	fi
}


########################################
# mysqlverify()
#	Parameters:
# 		$1 = base file name
#		$2 = app options to reproduce
# 		$3 = title
#       $4 = mysql validation query
########################################
mysqlverify () {
# Generate the mysql commands needed to validate...
cat >xxqq <<EOF
use mojosrv;
${4}
EOF
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $1 $3
	${CSVLOAD} $2 >>${LOGFILE} 2>&1
	mysql --no-defaults <xxqq >${1}

	if [ "${FORCEGOOD}" = "1" ]; then
		cp ${1} ${GOLD}/${1}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${1}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${1}.gold
			echo "Created a default $1.gold for you. Update this file with known-good output."
		fi
		UDIFFS=$(diff ${1} ${GOLD}/${1}.gold | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${CSVLOAD} ${2}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED...   if correct:  mv ${1} ${GOLD}/${1}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${CSVLOAD} ${2}" >> ${ERRFILE}
			echo "Differences in ${1} are as follows:" >> ${ERRFILE}
			diff ${GOLD}/${1}.gold ${1} >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${3}
			else
				if [ ${MANAGESERVER} -eq 1 ]; then
					echo "STOPPING MOJO SERVER"
					pkill mojosrv
				fi
				exit 1
			fi
		fi
	else
		echo
	fi
}

##########################################################################
# logcheck()
#   Compares log to log.gold
#   Date related fields are detected with a regular expression and changed
#   to "current time".  More filters may be needed depending on what goes
#   into the logfile.
#	Parameters:
#		none at this time
##########################################################################
logcheck() {
	echo -n "Test completed: " >> ${LOGFILE}
	date >> ${LOGFILE}
	if [ "${FORCEGOOD}" = "1" ]; then
		cp ${LOGFILE} ${GOLD}/${LOGFILE}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		echo -n "PHASE x: Log file check...  "
		if [ ! -f ${GOLD}/${LOGFILE}.gold -o ! -f ${LOGFILE} ]; then
			echo "Missing file -- Required files for this check: log.gold and log"
			failmsg
			exit 1
			# if [ "${ASKBEFOREEXIT}" = "1" ]; then
			# 	pause ${3}
			# else
			# 	if [ ${MANAGESERVER} -eq 1 ]; then
			# 		echo "STOPPING MOJO SERVER"
			# 		pkill mojosrv
			# 	fi
			# 	exit 1
			# fi
		fi
		declare -a out_filters=(
			's/^Date\/Time:.*/current time/'
			's/^Test completed:.*/current time/'
			's/(20[1-4][0-9]\/[0-1][0-9]\/[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] )(.*)/$2/'
			's/(20[1-4][0-9]\/[0-1][0-9]-[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] )(.*)/$2/'
			's/(20[1-4][0-9]-[0-1][0-9]-[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] )(.*)/$2/'
		)
		cp ${GOLD}/${LOGFILE}.gold ll.g
		cp log llog
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" ll.g > llx1; mv llx1 ll.g
			perl -pe "$f" llog > lly1; mv lly1 llog
		done
		UDIFFS=$(diff llog ll.g | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			echo "PASSED"
			passmsg
			rm -f ll.g llog
		else
			echo "FAILED:  differences are as follows:" >> ${ERRFILE}
			diff ll.g llog >> ${ERRFILE}
			echo >> ${ERRFILE}
			echo "If the new output is correct:  mv ${LOGFILE} ${GOLD}/${LOGFILE}.gold" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			exit 1
			# if [ "${ASKBEFOREEXIT}" = "1" ]; then
			# 	pause ${3}
			# else
			# 	if [ ${MANAGESERVER} -eq 1 ]; then
			# 		echo "STOPPING MOJO SERVER"
			# 		pkill mojosrv
			# 	fi
			# 	exit 1
			# fi
		fi
	else
		echo "FINISHED...  but did not check output"
	fi
	elapsedtime
}

##########################################################################
# genericlogcheck()
#   Compares the supplied file $1 to gold/$1.gold
#	Parameters:
# 		$1 = base file name
#		$2 = app options to reproduce
# 		$3 = title
##########################################################################
genericlogcheck() {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $1 $3
	if [ "${FORCEGOOD}" = "1" ]; then
		cp ${1} ${GOLD}/${1}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${1}.gold -o ! -f ${1} ]; then
			echo "Missing file -- Required files for this check: ${1} and ${GOLD}/${1}.gold"
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${3}
			else
				if [ ${MANAGESERVER} -eq 1 ]; then
					echo "STOPPING MOJO SERVER"
					pkill mojosrv
				fi
				exit 1
			fi
		fi
		UDIFFS=$(diff ${1} gold/${1}.gold | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			echo "PASSED"
		else
			echo "FAILED:  differences are as follows:" >> ${ERRFILE}
			diff gold/${1}.gold ${1} >> ${ERRFILE}
			echo >> ${ERRFILE}
			echo "If the new output is correct:  mv ${1} ${GOLD}/${1}.gold" >> ${ERRFILE}
			cat ${ERRFILE}
			exit 1
		fi
	else
		echo
	fi
}

#########################################################
# startMojoServer()
#	Kills any currently running instances of the server
#   then starts it up again.  The port is set to the
#   default port of 8270.  If you set MOJOPORT prior
#   to including base.sh to override the port number
#########################################################
startMojoServer () {
	if [ ${MANAGESERVER} -eq 1 ]; then
		stopMojoServer
		${MOJOBIN}/mojosrv -p ${MOJOPORT} > ${MOJOBIN}/mojolog 2>&1 &
		sleep 1
		rm -f mlog
		ln -s ${MOJOBIN}/mojolog mlog
	fi
}

#########################################################
# stopMojoServer()
#	Kills any currently running instances of the server
#########################################################
stopMojoServer () {
	if [ ${MANAGESERVER} -eq 1 ]; then
		killall mojosrv > /dev/null 2>&1
		sleep 1
	fi
}

########################################
# dojsonPOST()
#   Simulate a POST command to the server and use
#   the supplied file name as the json data
#	Parameters:
# 		$1 = url
#       $2 = json file
# 		$3 = base file name
#		$4 = title
########################################
dojsonPOST () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $3 $4
	CMD="curl -s -X POST ${1} -H \"Content-Type: application/json\" -d @${2}"
	${CMD} >rawcmdout; cat rawcmdout | python -m json.tool >${3} 2>>${LOGFILE}
	incStep

	if [ "${FORCEGOOD}" = "1" ]; then
		cp ${3} ${GOLD}/${3}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${3}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${3}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi

		#--------------------------------------------------------------------
		# The actual data has timestamp information that changes every run.
		# The timestamp can be filtered out for purposes of testing whether
		# or not the web service could be called and can return the expected
		# data.
		#--------------------------------------------------------------------
		declare -a out_filters=(
			's/(^[ \t]+"LastModTime":).*/$1 TIMESTAMP/'
		)
		cp gold/${3}.gold qqx
		cp ${3} qqy
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" qqx > qqx1; mv qqx1 qqx
			perl -pe "$f" qqy > qqy1; mv qqy1 qqy
		done

		UDIFFS=$(diff qqx qqy | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${CMD}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED..." >> ${ERRFILE}
			echo "Differences in ${3} are as follows:" >> ${ERRFILE}
			diff qqx qqy >> ${ERRFILE}
			echo "If correct:  mv ${3} ${GOLD}/${3}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${CMD}" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${3}
			else
				if [ ${MANAGESERVER} -eq 1 ]; then
					echo "STOPPING MOJO SERVER"
					pkill mojosrv
				fi
				exit 1
			fi
		fi
	else
		echo
	fi
	rm -f qqx qqy
}

########################################
# dojsonAwsPOST()
#   Simulate a POST command to the server and use
#   the supplied file name as the json data
#	Parameters:
# 		$1 = url
#       $2 = json file
# 		$3 = base file name
#		$4 = title
#		$5 = AWS SNS message type
########################################
dojsonAwsPOST () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $3 $4
	CMD="curl -s -X POST ${1} -H \"X-Amz-Sns-Message-Type: ${5}\" -d @${2}"
	HDR="X-Amz-Sns-Message-Type: ${5}"
	curl -s -X POST ${1} -H "${HDR}" -d @${2} | python -m json.tool >${3} 2>>${LOGFILE}
	incStep

	if [ "${FORCEGOOD}" = "1" ]; then
		cp ${3} ${GOLD}/${3}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${3}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${3}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi

		#--------------------------------------------------------------------
		# The actual data has timestamp information that changes every run.
		# The timestamp can be filtered out for purposes of testing whether
		# or not the web service could be called and can return the expected
		# data.
		#--------------------------------------------------------------------
		declare -a out_filters=(
			's/(^[ \t]+"LastModTime":).*/$1 TIMESTAMP/'
		)
		cp gold/${3}.gold qqx
		cp ${3} qqy
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" qqx > qqx1; mv qqx1 qqx
			perl -pe "$f" qqy > qqy1; mv qqy1 qqy
		done

		UDIFFS=$(diff qqx qqy | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${CMD}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED..." >> ${ERRFILE}
			echo "Differences in ${3} are as follows:" >> ${ERRFILE}
			diff qqx qqy >> ${ERRFILE}
			echo "If correct:  mv ${3} ${GOLD}/${3}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${CMD}" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${3}
			else
				if [ ${MANAGESERVER} -eq 1 ]; then
					echo "STOPPING MOJO SERVER"
					pkill mojosrv
				fi
				exit 1
			fi
		fi
	else
		echo
	fi
	rm -f qqx qqy
}

########################################
# dojsonGET()
#   Simulate a GET command to the server and use
#   the supplied file name as the json data
#	Parameters:
# 		$1 = url
# 		$2 = base file name
#		$3 = title
########################################
dojsonGET () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} ${2} ${3}
	CMD="curl -s \"${1}\""
	curl -s "${1}" | python -m json.tool >${2} 2>>${LOGFILE}
	incStep

	if [ "${FORCEGOOD}" = "1" ]; then
		cp ${2} ${GOLD}/${2}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${2}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${2}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi

		#--------------------------------------------------------------------
		# The actual data has timestamp information that changes every run.
		# The timestamp can be filtered out for purposes of testing whether
		# or not the web service could be called and can return the expected
		# data.
		#--------------------------------------------------------------------
		declare -a out_filters=(
			's/(^[ \t]+"LastModTime":).*/$1 TIMESTAMP/'
		)
		cp gold/${2}.gold qqx
		cp ${2} qqy
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" qqx > qqx1; mv qqx1 qqx
			perl -pe "$f" qqy > qqy1; mv qqy1 qqy
		done

		UDIFFS=$(diff qqx qqy | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${CMD}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED..." >> ${ERRFILE}
			echo "Differences in ${2} are as follows:" >> ${ERRFILE}
			diff qqx qqy >> ${ERRFILE}
			echo "If correct:  mv ${2} ${GOLD}/${2}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${CMD}" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${3}
			else
				if [ ${MANAGESERVER} -eq 1 ]; then
					echo "STOPPING MOJO SERVER"
					pkill mojosrv
				fi
				exit 1
			fi
		fi
	else
		echo
	fi
	rm -f qqx qqy
}
########################################
# dojsonGET()
#   Simulate a GET command to the server and use
#   the supplied file name as the json data
#	Parameters:
# 		$1 = url
# 		$2 = base file name
#		$3 = title
########################################
doHtmlGET () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} ${2} ${3}
	CMD="curl -s \"${1}\""
	curl -s "${1}" >${2} 2>>${LOGFILE}
	incStep

	if [ "${FORCEGOOD}" = "1" ]; then
		cp ${2} ${GOLD}/${2}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${2}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${2}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi

		#--------------------------------------------------------------------
		# The actual data has timestamp information that changes every run.
		# The timestamp can be filtered out for purposes of testing whether
		# or not the web service could be called and can return the expected
		# data.
		#--------------------------------------------------------------------
		declare -a out_filters=(
			's/(^[ \t]+"LastModTime":).*/$1 TIMESTAMP/'
		)
		cp gold/${2}.gold qqx
		cp ${2} qqy
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" qqx > qqx1; mv qqx1 qqx
			perl -pe "$f" qqy > qqy1; mv qqy1 qqy
		done

		UDIFFS=$(diff qqx qqy | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${CMD}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED..." >> ${ERRFILE}
			echo "Differences in ${2} are as follows:" >> ${ERRFILE}
			diff qqx qqy >> ${ERRFILE}
			echo "If correct:  mv ${2} ${GOLD}/${2}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${CMD}" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${3}
			else
				if [ ${MANAGESERVER} -eq 1 ]; then
					echo "STOPPING MOJO SERVER"
					pkill mojosrv
				fi
				exit 1
			fi
		fi
	else
		echo
	fi
	rm -f qqx qqy
}

#--------------------------------------------------------------------------
#  Handle command line options...
#--------------------------------------------------------------------------
tdir
while getopts "acfmornt:R:" o; do
	echo "o = ${o}"
	case "${o}" in
		a)	ASKBEFOREEXIT=1
			echo "WILL ASK BEFORE EXITING ON ERROR"
			;;
		c | C)
			SHOWCOMMAND=1
			echo "SHOWCOMMAND"
			;;
		r | R)
			doReport
			exit 0
			;;
		f)  SKIPCOMPARE=1
			echo "SKIPPING COMPARES..."
			;;
		m)  MANAGESERVER=0
			echo "SKIPPING SERVER MGMT CMDS..."
			;;
		n)	CREATENEWDB=0
			echo "DATA WILL BE ADDED TO CURRENT DB"
			;;
		o)	FORCEGOOD=1
			echo "OUTPUT OF THIS RUN IS SAVED AS *.GOLD"
			;;
		t) SINGLETEST="${OPTARG}"
			echo "SINGLETEST set to ${SINGLETEST}"
			;;
		*) 	usage
			exit 1
			;;
	esac
done
shift $((OPTIND-1))

if [ ! -f ${TREPORT} ]; then touch ${TREPORT}; fi
rm -f ${ERRFILE}
echo    "Test Name:    ${TESTNAME}" > ${LOGFILE}
echo    "Test Purpose: ${TESTSUMMARY}" >> ${LOGFILE}
echo -n "Date/Time:    " >>${LOGFILE}
date >> ${LOGFILE}
echo >>${LOGFILE}

echo -n "Create new database... " >> ${LOGFILE} 2>&1
if [ ${CREATENEWDB} -eq 1 ]; then
	${MOJOBIN}/mojonewdb
fi
if [ $? -eq 0 ]; then
	echo " successful" >> ${LOGFILE} 2>&1
else
	echo " ERROR" >> ${LOGFILE} 2>&1
	echo "Failed to create new database" > ${ERRFILE}
	cat ${ERRFILE}
	failmsg
	exit 1
fi
