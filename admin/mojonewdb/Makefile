TOP=../..


NEWDBPROG="mojonewdb"

mojonewdb: *.sh

clean:
	rm -f ${NEWDBPROG} fail
	@echo "*** Clean completed in admin/${NEWDBPROG} ***"

install: mojonewdb
	@echo "*** install completed in admin/${NEWDBPROG} ***"

package: mojonewdb
	touch fail
	cp ${NEWDBPROG} ${TOP}/tmp/mojo/
	@echo "*** Packaging completed in admin/${NEWDBPROG} ***"
	rm -f fail

test:
	@echo "*** Testing completed in admin/${NEWDBPROG} ***"
