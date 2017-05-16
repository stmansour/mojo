#!/bin/bash
CHECKINGPERIOD=30
TOTAL_TIME=0
COUNT=0
AVERAGE=0
LOGFILE=loadtime.csv

while [ 1=1 ];
do
	T=$(date "+%Y-%m-%d,%H:%M:%S")
	TIME=$(curl https://faatoday.com/ -o /dev/null -s -w %{time_total})
	TOTAL_TIME=$(echo "scale=5; ${TOTAL_TIME} + ${TIME}" |bc)
	COUNT=$((COUNT+1))
	AVERAGE=$(echo "scale=5; ${TOTAL_TIME} / ${COUNT}" |bc)
	DAT="$T,${TIME},${AVERAGE},${TOTAL_TIME}"

    #---------------------------------------------------------------------------
    # Touch the logfile, so we know that this process is active.
    # The timestamp on ${LOGFILE} shows when the process was last
    # checked.
    # Wait for a bit, then do it all again...
    #---------------------------------------------------------------------------
    echo "$DAT" >> ${LOGFILE}
    sleep ${CHECKINGPERIOD}
done
