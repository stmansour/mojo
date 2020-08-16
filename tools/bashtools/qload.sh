#!/bin/bash

BIN=../../tmp/mojo
F="music.csv"
TMP="xxqqmmjj"

${BIN}/mojonewdb
${BIN}/mojocsv -g smanmusic -cg -f "${TMP}"

