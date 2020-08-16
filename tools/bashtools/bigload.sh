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

echo "done"
