#!/bin/bash
# description: activation script to start/stop Accord RentRoll
#
# processname: mojo


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

stopwatchdog() {
	# make sure we can find it
    pidline=$(ps -ef | grep pbwatchdog | grep -v grep)
    if [ "${pidline}" != "" ]; then
        lines=$(echo "${pidline}" | wc -l)
        if [ $lines -gt "0" ]; then
            pid=$(echo "${pidline}" | awk '{print $2}')
            $(kill $pid)
        fi          
    fi      
}

makeProdNode() {
	${GETFILE} accord/db/config.json
}

#--------------------------------------------------------------
#  For QA, Sandbox, and Production nodes, go through the
#  laundry list of details...
#  1. Set up permissions for the database on QA and Sandbox nodes
#  2. Install a database with some data for testing
#  3. For PDF printing, install wkhtmltopdf
#--------------------------------------------------------------
setupAppNode() {
	./mojonewdb

	echo "Done."
}

start() {
	# Create a database if this is a localhost instance  
	if [ ${IAM} == "root" ]; then
		x=$(grep MojoDbhost config.json | grep localhost | wc -l)
		if (( x == 1 )); then
			setupAppNode
		fi
	fi

	if [ ${IAM} == "root" ]; then
		chown -R ec2-user *
		# chmod u+s ${PROGNAME} pbwatchdog
		if [ $(uname) == "Linux" -a ! -f "/etc/init.d/${PROGNAME}" ]; then
			cp ./activate.sh /etc/init.d/${PROGNAME}
			chkconfig --add ${PROGNAME}
		fi
	fi

	if [ ! -d "./js" ]; then
		${GETFILE} jenkins-snapshot/mojo/latest/js.tar.gz
		tar xzvf js.tar.gz
	fi
	if [ ! -d "./html/images" ]; then
		${GETFILE} jenkins-snapshot/mojo/latest/images.tar.gz
		tar xzvf images.tar.gz
	fi
	if [ ! -d "./html/fa" ]; then
		${GETFILE} jenkins-snapshot/mojo/latest/fa.tar.gz
		tar xzvf fa.tar.gz
	fi

	./${SERVERNAME} >log.out 2>&1 &

	# make sure it can survive a reboot...
	if [ ${IAM} == "root" ]; then
		if [ ! -d /var/run/${SERVERNAME} ]; then
			mkdir /var/run/${SERVERNAME}
		fi
		echo $! >/var/run/${SERVERNAME}/${SERVERNAME}.pid
		touch /var/lock/${SERVERNAME}
	fi

	# give ${SERVERNAME} a few seconds to start up before initiating the watchdog
	# sleep 5
	# if [ "${STARTPBONLY}" -ne "1" ]; then
	# 	stopwatchdog
	# 	./pbwatchdog ${WATCHDOGOPTS} >pbwatchdogstartup.out 2>&1 &
	# fi
}

stop() {
	# stopwatchdog
	killall -9 ${SERVERNAME}
	if [ ${IAM} == "root" ]; then
		sleep 2
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
  #   "images")
		# updateImages
		# echo "Images updated"
		# ;;
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
		# ST=$(curl -s http://${HOST}:${PORT}/status/)
		# echo "${ST}"
		echo "OK"
		exit 0
		;;
	# "status")
	# 	status
	# 	;;
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
	"updatedb")
		setupAppNode
		;;
	*)
		echo "Unrecognized command: $arg"
		usage
		exit 1
		;;
    esac
done
