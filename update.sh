#!/bin/bash
DEBUG=0

usage() {
    cat <<ZZEOF
ACCORD - update application script
Usage:   update.sh [OPTIONS]

OPTIONS:
-n  do not run startup after new version installed

Examples:
Normal usage
        bash$  sudo ./update.sh

Install but don't run startup script:
    bash$  ./pbrestore.sh -n

ZZEOF
        exit 0
}


#############################################################################
# decho
#   Description:
#   	Use this function like echo. If DEBUG is 1 then it will echo the
#		the output to the terminal. Otherwise it will just return without
#       echoing anything.
#
#   Params:
#		The string to echo
#
#	Returns:
#
#############################################################################
function decho {
	if (( ${DEBUG} == 1 )); then
		echo
		echo "${1}"
	fi
}

#############################################################################
# readConfig
#   Description:
#       Read the config.json file from the directory containing this script
#       to set some of the key values needed to access the artifactory repo.
#
#       Upon returning, URLBASE will end with the character "/".
#
#   Params:
#       none
#
#   Returns:
#       sets variables:  APIKEY, REPOUSER, and URLBASE
#
#############################################################################
readConfig() {
    RELDIR=$(cd `dirname "${BASH_SOURCE[0]}"` && pwd)
    CONF="${RELDIR}/config.json"
    REPOUSER=$(grep RepoUser ${CONF} | awk '{print $2;}' | sed -e 's/[,"]//g')
    APIKEY=$(grep RepoPass ${CONF} | awk '{print $2;}' | sed -e 's/[,"]//g')
    URLBASE=$(grep RepoURL ${CONF} | awk '{print $2;}' | sed -e 's/[,"]//g')

    # add a trailing / if it does not have one...
    if [ "${URLBASE: -1}" != "/" ]; then
        URLBASE="${URLBASE}/"
    fi

    pushd . >/dev/null
    if [ -d ~ec2-user ]; then
        cd ~ec2-user
        EC2USERHOME=$(pwd)
    else
        EC2USERHOME=$(pwd)
        echo "*** WARNING ***  home directory for ec2-user was not found"
        echo "                 will use ${EC2USERHOME} instead"
    fi
    popd >/dev/null

    decho "RELDIR = ${RELDIR}"
    decho "REPOUSER = ${REPOUSER}"
    decho "APIKEY = ${APIKEY}"
    decho "URLBASE = ${URLBASE}"
    decho "EC2USERHOME = ${EC2USERHOME}"
}

#############################################################################
# configure
#   Description:
#       The only configuration needed is the jfrog cli environment. Just
#       make sure we have it in the path. If it is not present, then
#       get it.
#
#   Params:
#       none
#
#   Returns:
#       nothing
#
#############################################################################
configure() {
    #---------------------------------------------------------
    # the user's bin directory is not created by default...
    #---------------------------------------------------------
    if [ ! -d ${EC2USERHOME}/bin ]; then
        mkdir ${EC2USERHOME}/bin
    fi

    #---------------------------------------------------------
    # now make sure that we have jfrog...
    #---------------------------------------------------------
    if [ ! -f ${EC2USERHOME}/bin/jfrog ]; then
        curl -s -u "${REPOUSER}:${APIKEY}" ${URLBASE}accord/tools/jfrog > ${EC2USERHOME}/bin/jfrog
        chown ec2-user:ec2-user ${EC2USERHOME}/bin/jfrog
        chmod +x ${EC2USERHOME}/bin/jfrog
    fi
    if [ ! -d ${EC2USERHOME}/.jfrog ]; then
        curl -s -u "${REPOUSER}:${APIKEY}" ${URLBASE}accord/tools/jfrogconf.tar > ${EC2USERHOME}/jfrogconf.tar
        pushd ${EC2USERHOME}
        tar xvf jfrogconf.tar
        rm jfrogconf.tar
        chown ec2-user:ec2-user ${EC2USERHOME}/bin/jfrog
        popd
    fi
    if [ ! -d ~root/.jfrog ]; then
        curl -s -u "${REPOUSER}:${APIKEY}" ${URLBASE}accord/tools/jfrogconf.tar > ~root/jfrogconf.tar
        pushd ~root
        tar xvf jfrogconf.tar
        rm jfrogconf.tar
        popd
    fi
    JFROG="${EC2USERHOME}/bin/jfrog"
    decho "JFROG = ${JFROG}"
}


#############################################################################
# GetLatestProductRelease
#   Description:
#       The only configuration needed is the jfrog cli environment. Just
#       make sure we have it in the path. If it is not present, then
#       get it.
#
#   Params:
#       ${1} = base name of product (rentroll, phonebook, mojo, ...)
#
#   Returns:
#       nothing
#
#############################################################################
GetLatestRepoRelease() {
    decho "GetLatestRepoRelease: searching for ${1}"
    f=$(${JFROG} rt s "accord/air/release/*" | grep ${1} | awk '{print $2}' | sed 's/"//g')
    if [ "x${f}" = "x" ]; then
        echo "Latest release of ${1}:  *** ERROR *** no release found"
        exit 1
    fi
    echo "Latest release of ${1}: ${f}"
    t=$(basename ${f})
    cd ${RELDIR}
    cdir=$(pwd)
    decho "Current working directory: ${cdir}"
    decho "curl -s -u ${REPOUSER}:${APIKEY} ${URLBASE}${f} > ../${t}"
    curl -s -u "${REPOUSER}:${APIKEY}" ${URLBASE}${f} > ../${t}
    decho "After call to curl, directory contents ls .."
    tmpx=$(ls ../phonebook_*.tar.gz)
    echo "Downlowded: ${tmpx}"
}

#############################################################################
# RunActivation
#   Description:
#       Runs sudo ./activate.sh start
#
#   Params:
#       ${1} = base name of product (rentroll, phonebook, mojo, ...)
#
#   Returns:
#       nothing
#
#############################################################################
RunActivation() {
	echo -n "Invoking activation script: "
	stat=$(./activate.sh -b start)
	sleep 2
	status=$(./activate.sh ready)
	if [ "${status}" = "OK" ]; then
	    echo "Success!"
	    rm ../phonebook*.tar.gz
	else
	    echo "error:  status = ${status}"
	    echo "output from ./activate.sh -b start "
	    echo "${stat}"
	fi
}

STARTUP=1
while getopts "n" o; do
        echo "o = ${o}"
        case "${o}" in
                n)      STARTUP=0
                        echo "WILL NOT RUN STARTUP AFTER INSTALL"
                        ;;
                *)      usage
                        exit 1
                        ;;
        esac
done
shift $((OPTIND-1))

readConfig
configure

#----------------------------------------------
#  ensure that we're in the phonebook directory...
#----------------------------------------------

cd ${RELDIR}
dir=${PWD##*/}
if [ ${dir} != "phonebook" ]; then
    echo "This script must execute in the phonebook directory."
    echo "current directory is: ${dir}"
    exit 1
fi

user=$(whoami)
if [ ${user} != "root" ]; then
    echo "This script must execute as root.  Try sudo !!"
    exit 1
fi

echo -n "Shutting down phonebook server: "
if [ -f "activate.sh" ]; then
    $(./activate.sh stop) >/dev/null 2>&1
    sleep 6
    echo "OK"
else
    echo "*** WARNING:  activate.sh was not found! Using killall instead ***"
    killall phonebook >/dev/null 2>&1
fi

cd ..
echo "Distribution download to:  ${PWD}"
rm -f phonebook*.tar*
GetLatestRepoRelease "phonebook"

echo -n "Extracting: "
cd ${RELDIR}/..
tar xzf phonebook*.tar.gz
chown -R ec2-user:ec2-user phonebook
cd ${RELDIR}
echo "done"

chmod u+s phonebook pbwatchdog

if [ "${STARTUP}" = "1" ]; then
        RunActivation
fi
