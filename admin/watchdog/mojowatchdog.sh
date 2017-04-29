#!/bin/bash
PORT=8275
CHECKINGPERIOD=10                  # seconds
LOGFILE="mojowatchdog.log"

#----------------------------------------------------
#  Main loop:   
#  Ping mojosrv on localhost:8275.
#  If we don't hear back, then restart
#----------------------------------------------------
while [ 1=1 ];
do
	R=$(curl -s http://localhost:${PORT}/v1/ping | grep "Accord Mojo" | wc -l)
	if [ 0 = ${R} ]; then
		echo -n "Ping to mojosrv failed at " >> ${LOGFILE}
		date >>  ${LOGFILE}
		echo -n "Restart..." >> ${LOGFILE}
		pkill mojosrv
		./activate.sh start
	fi

    #---------------------------------------------------------------------------
    # Touch the logfile, so we know that this process is active.
    # The timestamp on ${LOGFILE} shows when the process was last
    # checked.
    # Wait for a bit, then do it all again...
    #---------------------------------------------------------------------------
    touch ${LOGFILE}
    sleep ${CHECKINGPERIOD}
done
