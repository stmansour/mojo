#!/bin/bash
# description: activation script to start/stop Accord RentRoll
#
# processname: mojo

OS=$(uname)
HOST=localhost
PROGNAME="mojo"
PORT=8275
WATCHDOGOPTS=""
GETFILE="/usr/local/accord/bin/getfile.sh"
RENTROLLHOME="/home/ec2-user/apps/${PROGNAME}"
DATABASENAME="${PROGNAME}"
DBUSER="ec2-user"
SERVERNAME="mojosrv"
IAM=$(whoami)
TRACE=0


usage() {
    cat <<ZZEOF
Mojo activation script.
Usage:   activate.sh [OPTIONS] CMD

This is the Accord Mojo activation script. It is designed to work in two environments.
First, it works with Plum - Accord's test environment automation infrastructure
Second, it can work as a service script in /etc/init.d

OPTIONS:
-p port      (default is 8275)
-h hostname  (default is localhost)
-N dbname    (default is ${PROGNAME})
-T           (use this option to indicate testing rather than production)

CMD is one of: start | stop | status | restart | ready | reload | condrestart | makeprod



Examples:
Command to start ${PROGNAME}:
	bash$  activate.sh start

Command to stop ${PROGNAME}:
	bash$  activate.sh Stop

Command to see if ${PROGNAME} is ready for commands... the response
will be "OK" if it is ready, or something else if there are problems:

    bash$  activate.sh ready
    OK
ZZEOF
}

makeProdNode() {
	${GETFILE} accord/db/confprod.json
	cp confprod.json config.json
}

#--------------------------------------------------------------
#  For QA, Sandbox, and Production nodes, go through the
#  laundry list of details...
#  1. Set up permissions for the database on QA and Sandbox nodes
#  2. Install a database with some data for testing
#  3. For PDF printing, install wkhtmltopdf
#--------------------------------------------------------------
makeDevNode() {
	${GETFILE} accord/db/confdev.json
	cp confdev.json config.json
	./mojonewdb
	echo "Done."
}

Trace () {
    if (( TRACE == 1 )); then
        echo $1
    fi
}

start() {
	# Create a database if this is a localhost instance
    Trace "Entered start"
	if [ ! -f "config.json" ]; then
		echo "config.json not found, setting up as development node"
		makeDevNode
	fi

    Trace "[2]"

	if [ ${IAM} == "root" ]; then
		if [ ! -f "mojo.log" ]; then
            Trace "[3]"
			touch mojo.log
			touch mojowatchdog.log
		fi
        Trace "[4]"
		chown -R ec2-user:ec2-user *
		# chmod u+s ${PROGNAME} pbwatchdog
        Trace "[5]"
		if [ $(uname) == "Linux" -a ! -f "/etc/init.d/${PROGNAME}" ]; then
            Trace "[6]"
			cp ./activate.sh /etc/init.d/${PROGNAME}
			chkconfig --add ${PROGNAME}
		fi
	fi
    Trace "[7]"

	# if [ ! -d "./js" ]; then
	# 	${GETFILE} jenkins-snapshot/mojo/latest/js.tar.gz
	# 	tar xzvf js.tar.gz
	# fi
	# if [ ! -d "./html/images" ]; then
	# 	${GETFILE} jenkins-snapshot/mojo/latest/images.tar.gz
	# 	tar xzvf images.tar.gz
	# fi
	# if [ ! -d "./html/fa" ]; then
	# 	${GETFILE} jenkins-snapshot/mojo/latest/fa.tar.gz
	# 	tar xzvf fa.tar.gz
	# fi

	x=$(pgrep "${SERVERNAME}")
    Trace "[8]"
	if [ "${X}x" == "x" ]; then
        Trace "[9]"
		./${SERVERNAME} >log.out 2>&1 &
        Trace "[10]"
	fi

    Trace "[11]"
	# make sure it can survive a reboot...
	if [ ${IAM} == "root" ]; then
        Trace "[12]"
		if [ ! -d /var/run/${SERVERNAME} ]; then
            Trace "[13]"
			mkdir /var/run/${SERVERNAME}
		fi
        Trace "[14]"
		echo $! >/var/run/${SERVERNAME}/${SERVERNAME}.pid
		touch /var/lock/${SERVERNAME}
	fi

    Trace "[15]"
	# give ${SERVERNAME} a few seconds to start up before initiating the watchdog
	sleep 1

    Trace "[16]"
	#---------------------------------------------------
	# If the watchdog is NOT running, then start it...
	#---------------------------------------------------
	W=$(ps -ef | grep "mojowatchdog" | grep "bash" | wc -l)
    Trace "[17]"
	if [ ${W} == 0 ]; then
        Trace "[18]"
		./mojowatchdog &
        Trace "[19]"
	fi
    Trace "[20]"
}

stop() {
	#---------------------------------------------------
	# stop watchdog first
	#---------------------------------------------------
	W=$(ps -ef | grep "mojowatchdog" | grep "bash" | wc -l)
	if [ ${W} == 1 ]; then
		case "${OS}" in
		"Darwin")
			pid=$(ps -ef | grep mojowatchdog | grep "bash" | sed -e 's/[ \t]*[0-9][0-9]*[ \t][ \t]*\([0-9][0-9]*\)[ \t].*/\1/')
			;;
		"Linux")
			pid=$(ps -ef | grep mojowatchdog | grep "bash" | sed -e 's/[^ \t]*[ \t][ \t]*\([0-9][0-9]*\)[ \t].*/\1/')
			;;
		"*")
			echo "Unsupported Operating System"
			exit 1
			;;
		esac
		kill ${pid}
	fi

	#---------------------------------------------------
	# now stop the server
	#---------------------------------------------------
	pkill ${SERVERNAME}
	sleep 1
	X=$(pgrep ${SERVERNAME})
	if [ "x${X}" != "x" ]; then
		killall -9 ${SERVERNAME}
	fi

	if [ ${IAM} == "root" ]; then
		sleep 1
		rm -f /var/run/${SERVERNAME}/${SERVERNAME}.pid /var/lock/${SERVERNAME}
	fi
}

status() {
	ST=$(curl -s http://${HOST}:${PORT}/v1/ping/ | wc -c)
	case "${ST}" in
	"33")
		exit 0
		;;
	"0")
		# ${SERVERNAME} is not responsive or not running.  Exit status as described in
		# http://refspecs.linuxbase.org/LSB_3.1.0/LSB-Core-generic/LSB-Core-generic/iniscrptact.html
		if [ ${IAM} == "root" -a -f /var/run/${SERVERNAME}/${SERVERNAME}.pid ]; then
			exit 1
		fi
		if [ ${IAM} == "root" -a -f /var/lock/${SERVERNAME} ]; then
			exit 2
		fi
		exit 3
		;;
	"*") echo "Not sure what state it's in. Response had ${ST} characters, expected 33."
	esac
}

reload() {
	ST=$(curl -s http://${HOST}:${PORT}/v1/ping)
}

restart() {
	stop
	sleep 3
	start
}

Trace "TRACE IS TURNED ON"

while getopts ":p:qih:N:Tb" o; do
    case "${o}" in
       b)
            WATCHDOGOPTS="-b"
	    	# echo "WATCHDOGOPTS set to: ${WATCHDOGOPTS}"
            ;;
       h)
            HOST=${OPTARG}
            echo "HOST set to: ${HOST}"
            ;;
        N)
            DATABASENAME=${OPTARG}
            # echo "DATABASENAME set to: ${DATABASENAME}"
            ;;
        p)
            PORT=${OPTARG}
	    	# echo "PORT set to: ${PORT}"
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

# cd "${RENTROLLHOME}"
PBPATH=$(cd `dirname "${BASH_SOURCE[0]}"` && pwd)
cd ${PBPATH}

for arg do
	# echo '--> '"\`$arg'"
	cmd=$(echo ${arg}|tr "[:upper:]" "[:lower:]")
    case "$cmd" in
	"start")
		start
		echo "OK"
		exit 0
		;;
	"stop")
		stop
		echo "OK"
		exit 0
		;;
	"ready")
		R=$(curl -s http://localhost:${PORT}/v1/ping | grep "Accord Mojo" | wc -l)
		if [ 1 = ${R} ]; then
			echo "OK"
		else
			echo "No response to ping"
		fi
		exit 0
		;;
	"restart")
		restart
		echo "OK"
		exit 0
		;;
	"reload")
		reload
		exit 0
		;;
	"condrestart")
		if [ -f /var/lock/mojosrv ] ; then
			restart
		fi
		;;
	"makeprod")
		makeProdNode
		;;
	"makedev")
		makeDevNode
		;;
	*)
		echo "Unrecognized command: $arg"
		usage
		exit 1
		;;
    esac
done
