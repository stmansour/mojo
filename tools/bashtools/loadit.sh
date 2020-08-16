#!/bin/bash
BIN=../../tmp/mojo
F="music.csv"
TMP="xxqqmmjj"

#-----------------------------------------------
# Start a new db with the database in mojo.sql
#-----------------------------------------------
echo "Create new database from mojo.sql"
${BIN}/mojonewdb
mysql --no-defaults mojo < mojo.sql

#------------------------------------------
# Extract just the smanmusic entries...
#------------------------------------------
echo "Extract smanmusic entries"
${BIN}/mojoexport -group smanmusic > ${F}

#------------------------------------------
# Recreate the db with just those entries
#------------------------------------------
echo "Create newdb with smanmusic entries only"

# we need to create the map header line -- same as first line but remove
# the status column...
MAP=$(head -1 "${F}" | sed 's/^[^,][^,]*,/,/')
echo "MAP = ${MAP}"

rm -f ${TMP}

COUNT=0
while IFS= read -r line
do
    echo "${line}" >> "${TMP}"
    ((COUNT++))
    if [ "${COUNT}" = "1" ]; then
        echo "${MAP}" >> ${TMP}
    fi
done < "${F}"

${BIN}/mojonewdb
${BIN}/mojocsv -g smanmusic -cg -f "${TMP}"

echo "done"
