PASS=AP3wHZhcQQCvkC4GVCCZzPcqe3L
ART=http://ec2-52-91-201-195.compute-1.amazonaws.com/artifactory
GETFILE="/usr/local/accord/bin/getfile.sh"
USR=accord
PRODUCT=mojo

EXTERNAL_HOST_NAME=$(curl -s http://169.254.169.254/latest/meta-data/public-hostname)
#${EXTERNAL_HOST_NAME:?"Need to set EXTERNAL_HOST_NAME non-empty"}

#--------------------------------------------------------------
#  Routine to download files from Artifactory
#--------------------------------------------------------------
artf_get() {
    echo "Downloading $1/$2"
    wget -O "$2" --user=$USR --password=$PASS ${ART}/"$1"/"$2"
}

loadAccordTools() {
    #--------------------------------------------------------------
    #  Let's get our tools in place...
    #--------------------------------------------------------------
    artf_get ext-tools/utils accord-linux.tar.gz
    echo "Installing /usr/local/accord" >>${LOGFILE}
    cd /usr/local
    tar xzf ~ec2-user/accord-linux.tar.gz
    chown -R ec2-user:ec2-user accord
    cd ~ec2-user/
}

#----------------------------------------------
#  ensure that we're in the ${PRODUCT} directory...
#----------------------------------------------
dir=${PWD##*/}
if [ ${dir} != "${PRODUCT}" ]; then
    echo "This script must execute in the ${PRODUCT} directory."
    exit 1
fi

user=$(whoami)
if [ ${user} != "root" ]; then
    echo "This script must execute as root.  Try sudo !!"
    exit 1
fi

echo -n "Shutting down ${PRODUCT} server."; $(./activate.sh stop) >/dev/null 2>&1
echo -n "."
echo -n "."; 
echo -n "."; cd ..
echo
echo -n "Retrieving latest development snapshot of ${PRODUCT}..."
${GETFILE} jenkins-snapshot/${PRODUCT}/latest/${PRODUCT}.tar.gz
echo
echo -n "."; gunzip -f ${PRODUCT}.tar.gz
echo -n "."; tar xf ${PRODUCT}.tar
echo -n "."; chown -R ec2-user:ec2-user ${PRODUCT}
echo -n "."; cd ${PRODUCT}/
echo -n "."; echo -n "starting..."
echo -n "."; ./activate.sh start
echo -n "."; sleep 1
echo -n "."; status=$(./activate.sh ready)
echo
#  ./installman.sh >installman.log 2>&1
if [ "${status}" = "OK" ]; then
    echo "Activation successful"
else
    echo "Problems activating ${PRODUCT}.  Status = ${status}"
fi
# ${GETFILE} jenkins-snapshot/${PRODUCT}/latest/images.tar.gz
# tar xzvf rrimages.tar.gz
# ${GETFILE} jenkins-snapshot/${PRODUCT}/latest/js.tar.gz
# tar xzvf rrjs.tar.gz
# ${GETFILE} jenkins-snapshot/${PRODUCT}/latest/fa.tar.gz
# tar xzvf fa.tar.gz
