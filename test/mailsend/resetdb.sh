#!/bin/bash
if [ -f ./mojonewdb ]; then
    ./mojonewdb
elif [ -f ../../tmp/mojo/mojonewdb ]; then
    ../../tmp/mojo/mojonewdb
fi

./mailsend -setup
